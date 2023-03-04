package logger

import (
	"os"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logger *zerolog.Logger
var once sync.Once

func NewLogger(logFile string) *zerolog.Logger {
	once.Do(func() {
		output := zerolog.ConsoleWriter{Out: os.Stdout}
		if logFile != "" {
			file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Error().Msgf("Failed to open log file: %s", err)
			}
			output = zerolog.ConsoleWriter{Out: file}
		}

		//if first time, create logger.
		newLogger := zerolog.New(output).With().Timestamp().Logger()
		//save it to global var.
		logger = &newLogger
	})

	return logger
}
