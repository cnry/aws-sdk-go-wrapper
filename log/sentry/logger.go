/*
	for logging to sentry
	use like this
		import "github.com/cnry/aws-sdk-go-wrapper/"
		import _ "github.com/cnry/aws-sdk-go-wrapper/log/sentry"
*/

package sentry

import (
	LOG "github.com/cnry/aws-sdk-go-wrapper/log"
	Sentry "github.com/cnry/go-sentry-logger"
)

// override loggers in initialize
func init() {
	LOG.SetLogger(newSentryLogger())
}

type SentryLogger struct{}

func newSentryLogger() *SentryLogger {
	return &SentryLogger{}
}

func (l *SentryLogger) Fatal(label string, value LOG.Any) {
	data := Sentry.NewLogData(value, 3)
	data.Label = label
	Sentry.Fatal(data)
}

func (l *SentryLogger) Error(label string, value LOG.Any) {
	data := Sentry.NewLogData(value, 3)
	data.Label = label
	Sentry.Error(data)
}

func (l *SentryLogger) Warn(label string, value LOG.Any) {
	data := Sentry.NewLogData(value, 3)
	data.Label = label
	Sentry.Warn(data)
}

func (l *SentryLogger) Info(label string, value LOG.Any) {
	data := Sentry.NewLogData(value, 3)
	data.Label = label
	Sentry.Info(data)
}
