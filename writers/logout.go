package writers

import (
	"github.com/sirupsen/logrus"
	"github.com/tonydonlon/eventservice/api"
	"github.com/tonydonlon/eventservice/logger"
)

var log *logrus.Logger

func init() {
	log = logger.GetLogger()
}

// StdOutWriter implements events.EventWriter and writes to stdout
type StdOutWriter struct{}

func (s *StdOutWriter) Write(msg api.Event, sessionID string) error {
	log.WithFields(logrus.Fields{
		"name":  "StdOutWriter",
		"event": msg,
	}).Info("writing event")
	return nil
}
