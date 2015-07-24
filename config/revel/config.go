package revel

import (
	loader "github.com/cnry/revel-config-loader"

	"github.com/cnry/aws-sdk-go-wrapper/config"
)

const (
	awsConfigSectionName = "aws"
)

func init() {
	config.SetConfig(NewConfig())
}

// Config is struct for revel config
type Config struct{}

// NewConfig creates new Config for revel config
func NewConfig() *Config {
	return &Config{}
}

// GetConfigValue gets value from loaded revel configs
func (c *Config) GetConfigValue(section, key, df string) string {
	return loader.GetConfigValueDefault(config.AWSConfigFileName, awsConfigSectionName, section+"."+key, df)
}

// TODO: SetValues
func (c *Config) SetValues(data map[string]interface{}) {
	// TODO: implement
}

// TODO: LoadFile
func (c *Config) LoadFile(path string) error {
	// TODO: implement
	return nil
}
