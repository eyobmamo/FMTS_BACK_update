package utils

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger interface {
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
	Debugf(template string, args ...interface{})
	Sync() error
}

type zeroLogger struct {
	logger zerolog.Logger
}

func NewLogger() Logger {
	zerolog.TimeFieldFormat = time.RFC3339
	logger := zerolog.New(os.Stderr).
		With().
		Timestamp().
		Str("service", "go_starter").
		Logger()

	return &zeroLogger{logger}
}

func (l *zeroLogger) Infof(template string, args ...interface{}) {
	l.logger.Info().Msgf(template, args...)
}

func (l *zeroLogger) Warnf(template string, args ...interface{}) {
	l.logger.Warn().Msgf(template, args...)
}

func (l *zeroLogger) Errorf(template string, args ...interface{}) {
	l.logger.Error().Msgf(template, args...)
}

func (l *zeroLogger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatal().Msgf(template, args...)
	os.Exit(1)
}

func (l *zeroLogger) Debugf(template string, args ...interface{}) {
	l.logger.Debug().Msgf(template, args...)
}

func (l *zeroLogger) Sync() error {
	return nil
}
