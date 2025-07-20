package utils

import (
	"log"
	"os"
)

type StandardLogger struct {
	logger *log.Logger
}

type Logger interface {
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
	Debugf(template string, args ...interface{})
	Sync() error
}

func NewStandardLogger() *StandardLogger {
	return &StandardLogger{
		logger: log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
	}
}

func (l *StandardLogger) Infof(template string, args ...interface{}) {
	l.logger.Printf("[INFO] "+template, args...)
}

func (l *StandardLogger) Warnf(template string, args ...interface{}) {
	l.logger.Printf("[WARN] "+template, args...)
}

func (l *StandardLogger) Errorf(template string, args ...interface{}) {
	l.logger.Printf("[ERROR] "+template, args...)
}

func (l *StandardLogger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatalf("[FATAL] "+template, args...)
}

func (l *StandardLogger) Debugf(template string, args ...interface{}) {
	l.logger.Printf("[DEBUG] "+template, args...)
}

func (l *StandardLogger) Sync() error {
	// No-op for standard log, but required by interface
	return nil
}
