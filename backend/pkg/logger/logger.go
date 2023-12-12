package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

var ZL zerolog.Logger

func New(level string) {
	lvl := zerolog.InfoLevel
	switch level {
	case "debug":
		lvl = zerolog.DebugLevel
	case "warn":
		lvl = zerolog.WarnLevel
	case "error":
		lvl = zerolog.ErrorLevel
	}

	ZL = zerolog.New(
		zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		},
	).
		Level(lvl).
		With().
		Timestamp().
		Caller().
		Logger()
}
