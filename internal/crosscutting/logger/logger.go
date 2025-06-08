package logger

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

const (
	logLevelInfo  = "INFO"
	logLevelError = "ERROR"
)

type ILogger interface {
	Info(msg string)
	Error(msg string)
}

type logEntry struct {
	Level     string `json:"level"`
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
}

type StdLogger struct {
	logger *log.Logger
}

func NewLogger() ILogger {
	return &StdLogger{logger: log.New(os.Stdout, "", 0)}
}

func (l *StdLogger) log(level, msg string) {
	entry := logEntry{
		Level:     level,
		Timestamp: time.Now().Format(time.RFC3339),
		Message:   msg,
	}
	b, err := json.Marshal(entry)
	if err != nil {
		l.logger.Printf(`{"level":"ERROR","timestamp":"%s","message":"failed to marshal log entry: %v"}`, time.Now().Format(time.RFC3339), err)
		return
	}
	l.logger.Println(string(b))
}

func (l *StdLogger) Info(msg string) {
	l.log(logLevelInfo, msg)
}

func (l *StdLogger) Error(msg string) {
	l.log(logLevelError, msg)
}
