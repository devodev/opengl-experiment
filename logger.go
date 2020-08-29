package main

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

// ErrorType .
type ErrorType int

// ErrorType enum .
const (
	GLFWInitError ErrorType = iota + 1
	GLFWCreateWindowError
	GLInitError
	GLCompileShaderrError
)

// FatalError logs.
func (l *SimpleLogger) FatalError(e ErrorType, err error) int {
	errorMessages := map[ErrorType]string{
		GLFWInitError:         "could not initialize GLFW",
		GLFWCreateWindowError: "could not create window",
		GLInitError:           "could not initialize OpenGL",
		GLCompileShaderrError: "could not compile shader",
	}
	l.Errorf(errorMessages[e]+": %s", err)

	return int(e)
}
