package logger

import (
	"os"

	log "github.com/sirupsen/logrus"
)

// TODO logger proxy to mixin named logger plus correlationID

// GetLogger returns a logrus configured logger
func GetLogger() *log.Logger {
	checkLevel := func(level string) log.Level {
		for _, lvl := range log.AllLevels {
			if lvl.String() == level {
				return lvl
			}
		}
		log.Infof("Unknown LOG_LEVEL: '%s' using default 'info'", level)
		return log.InfoLevel
	}

	var logLevel = checkLevel(os.Getenv("LOG_LEVEL"))

	return &log.Logger{
		Formatter: &log.JSONFormatter{
			PrettyPrint: logLevel == log.DebugLevel,
		},
		Level: log.DebugLevel,
		Out:   os.Stdout,
	}
}
