package engine

import (
	"github.com/sirupsen/logrus"
)

// SimpleLogger .
type SimpleLogger struct {
	*logrus.Logger
}

// NewLogger .
func NewLogger() *SimpleLogger {
	logger := logrus.StandardLogger()
	logger.SetReportCaller(true)

	return &SimpleLogger{logger}
}
