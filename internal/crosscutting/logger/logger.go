package logger

import (
	"log"
)

const (
	logLevelInfo  = "INFO"
	logLevelError = "ERROR"
	colorGreen    = "\033[32m"
	colorRed      = "\033[31m"
	colorReset    = "\033[0m"
)

type ILogger interface {
	Info(msg string)
	Error(msg string)
}

type StdLogger struct{}

func (l *StdLogger) Info(msg string) {
	log.Printf("%s%s%s: %s", colorGreen, logLevelInfo, colorReset, msg)
}

func (l *StdLogger) Error(msg string) {
	log.Printf("%s%s%s: %s", colorRed, logLevelError, colorReset, msg)
}

func NewLogger() ILogger {
	return &StdLogger{}
}
