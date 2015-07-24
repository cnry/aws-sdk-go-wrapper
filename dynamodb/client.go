// DynamoDB Client

package dynamodb

import (
	"errors"
	"strings"

	SDK "github.com/aws/aws-sdk-go/service/dynamodb"

	"github.com/cnry/aws-sdk-go-wrapper/auth"
	"github.com/cnry/aws-sdk-go-wrapper/config"
	"github.com/cnry/aws-sdk-go-wrapper/log"
)

const (
	dynamodbConfigSectionName = "dynamodb"
	defaultRegion             = "us-east-1"
	defaultEndpoint           = "http://localhost:8000"
	defaultTablePrefix        = "dev_"
)

// wrapped struct for DynamoDB
type AmazonDynamoDB struct {
	TablePrefix string

	client      *SDK.DynamoDB
	tables      map[string]*DynamoTable
	writeTables map[string]bool
}

// Create new AmazonDynamoDB struct
func NewClient() *AmazonDynamoDB {
	region := config.GetConfigValue(dynamodbConfigSectionName, "region", "")
	endpoint := config.GetConfigValue(dynamodbConfigSectionName, "endpoint", "")
	conf := auth.NewConfig(region, endpoint)
	return newClient(conf)
}

// Create new AmazonDynamoDB struct
func NewClientWithKeys(k auth.Keys) *AmazonDynamoDB {
	conf := auth.NewConfigWithKeys(k)
	return newClient(conf)
}

// Create new AmazonDynamoDB struct
func newClient(conf auth.Config) *AmazonDynamoDB {
	d := &AmazonDynamoDB{}
	d.tables = make(map[string]*DynamoTable)
	d.writeTables = make(map[string]bool)
	d.TablePrefix = config.GetConfigValue(dynamodbConfigSectionName, "prefix", defaultTablePrefix)

	conf.SetDefault(defaultRegion, defaultEndpoint)
	awsConf := conf.Config
	d.client = SDK.New(awsConf)
	return d
}

// Create new DynamoDB table
func (d *AmazonDynamoDB) CreateTable(in *SDK.CreateTableInput) error {
	data, err := d.client.CreateTable(in)
	if err != nil {
		log.Error("[DynamoDB] Error on `CreateTable` operation, table="+*in.TableName, err)
		return err
	}
	log.Info("[DynamoDB] Complete CreateTable, table="+*in.TableName, data.TableDescription.TableStatus)
	return nil
}

// Delete DynamoDB table
func (d *AmazonDynamoDB) DeleteTable(name string) error {
	in := &SDK.DeleteTableInput{
		TableName: String(name),
	}
	data, err := d.client.DeleteTable(in)
	if err != nil {
		log.Error("[DynamoDB] Error on `DeleteTable` operation, table="+*in.TableName, err)
		return err
	}
	log.Info("[DynamoDB] Complete DeleteTable, table="+*in.TableName, data.TableDescription.TableStatus)
	return nil
}

// get infomation of the table
func (d *AmazonDynamoDB) DescribeTable(name string) (*SDK.TableDescription, error) {
	req, err := d.client.DescribeTable(&SDK.DescribeTableInput{
		TableName: String(name),
	})
	if err != nil {
		return nil, err
	}
	return req.Table, nil
}

// get the DynamoDB table
func (d *AmazonDynamoDB) GetTable(table string) (*DynamoTable, error) {
	tableName := d.TablePrefix + table

	// get the table from cache
	t, ok := d.tables[tableName]
	if ok {
		return t, nil
	}

	// get the table info from server
	desc, err := d.DescribeTable(tableName)
	if err != nil {
		return nil, err
	}
	t = &DynamoTable{
		db:      d,
		table:   desc,
		name:    tableName,
		indexes: make(map[string]*DynamoIndex),
	}
	for _, idx := range desc.LocalSecondaryIndexes {
		t.indexes[*idx.IndexName] = NewDynamoIndex(*idx.IndexName, indexTypeLSI, idx.KeySchema)
	}
	for _, idx := range desc.GlobalSecondaryIndexes {
		t.indexes[*idx.IndexName] = NewDynamoIndex(*idx.IndexName, indexTypeGSI, idx.KeySchema)
	}
	d.tables[tableName] = t
	return t, nil
}

// add the table to write spool
func (d *AmazonDynamoDB) addWriteTable(name string) {
	d.writeTables[name] = true
}

// remove the table from write spool
func (d *AmazonDynamoDB) removeWriteTable(name string) {
	d.writeTables[name] = false
}

// execute put operation for all tables in write spool
func (d *AmazonDynamoDB) PutAll() error {
	var errs []string
	for name, _ := range d.writeTables {
		err := d.tables[name].Put()
		if err != nil {
			errs = append(errs, err.Error())
			log.Error("[DynamoDB] Error on `Put` operation, table="+name, err.Error())
		}
		d.removeWriteTable(name)
	}
	if len(errs) != 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

// get the list of DynamoDB table
func (d *AmazonDynamoDB) ListTables() ([]*string, error) {
	res, err := d.client.ListTables(&SDK.ListTablesInput{})
	if err != nil {
		return make([]*string, 0, 0), err
	}
	return res.TableNames, nil
}
