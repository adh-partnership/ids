package signalr

import (
	"fmt"

	"github.com/adh-partnership/ids/backend/pkg/logger"
	"github.com/rs/zerolog"
)

type Logger struct{}

// Wrapper for zerolog
func (l *Logger) Log(keyvals ...interface{}) error {
	var log *zerolog.Event

	mapKeyvals := make(map[string]string)
	for i := 0; i < len(keyvals); i += 2 {
		mapKeyvals[fmt.Sprintf("%v", keyvals[i])] = fmt.Sprint(keyvals[i+1])
	}

	switch mapKeyvals["level"] {
	case "debug":
		log = logger.ZL.Debug()
	case "warn":
		log = logger.ZL.Warn()
	case "error":
		log = logger.ZL.Error()
	default:
		log = logger.ZL.Info()
	}

	for k, v := range mapKeyvals {
		if k != "level" {
			log = log.Str(k, v)
		}
	}

	log.Msg("SignalR")

	return nil
}
