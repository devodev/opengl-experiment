package main

import (
	"os"

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

// ErrorType .
type ErrorType int

// ErrorType enum .
const (
	GLFWInitError ErrorType = iota + 1
	GLInitError
	GLFWCreateWindowError
)

// FatalError logs.
func (l *SimpleLogger) FatalError(e ErrorType, err error) {
	errorMessages := map[ErrorType]string{
		GLFWInitError:         "could not initialize GLFW",
		GLInitError:           "could not initialize OpenGL",
		GLFWCreateWindowError: "could not create window",
	}
	l.Errorf(errorMessages[e]+": %s", err)
	os.Exit(int(e))
}
