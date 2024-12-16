package app

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

type TimestampFormatter struct{}

func NewZeroLog() zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	logger := zerolog.New(os.Stdout).With().Caller().Logger().Hook(&TimestampFormatter{})

	return logger
}

func (t *TimestampFormatter) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	e.Str("time", time.Now().Format("2006-01-02 15:04:05"))
}
