package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, err error, fields map[string]interface{})
	Fatal(msg string, err error, fields map[string]interface{})
	With(fields map[string]interface{}) Logger
}

type zeroLogger struct {
	logger zerolog.Logger
}

func NewLogger(serviceName, version string, debug bool) Logger {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.MessageFieldName = "message"
	zerolog.ErrorFieldName = "error"

	logLevel := zerolog.InfoLevel
	if debug {
		logLevel = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(logLevel)
	multi := zerolog.MultiLevelWriter(os.Stdout)

	if debug {
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
		multi = zerolog.MultiLevelWriter(consoleWriter)
	}

	logger := zerolog.
		New(multi).With().
		Timestamp().
		Str("service", serviceName).
		Str("version", version).
		Logger()

	return &zeroLogger{
		logger: logger,
	}
}

func (l *zeroLogger) Debug(msg string, fields map[string]interface{}) {
	l.logger.Debug().Fields(fields).Msg(msg)
}

func (l *zeroLogger) Info(msg string, fields map[string]interface{}) {
	l.logger.Info().Fields(fields).Msg(msg)
}

func (l *zeroLogger) Warn(msg string, fields map[string]interface{}) {
	l.logger.Warn().Fields(fields).Msg(msg)
}

func (l *zeroLogger) Error(msg string, err error, fields map[string]interface{}) {
	l.logger.Error().Err(err).Fields(fields).Msg(msg)
}

func (l *zeroLogger) Fatal(msg string, err error, fields map[string]interface{}) {
	l.logger.Fatal().Err(err).Fields(fields).Msg(msg)
}

func (l *zeroLogger) With(fields map[string]interface{}) Logger {
	return &zeroLogger{
		logger: l.logger.With().Fields(fields).Logger(),
	}
}
