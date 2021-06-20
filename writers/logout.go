package writers

import (
	"github.com/sirupsen/logrus"
	"github.com/tonydonlon/eventservice/logger"
	"github.com/tonydonlon/eventservice/streams"
)

var log *logrus.Logger

func init() {
	log = logger.GetLogger()
}

// StdOutWriter implements events.EventWriter and writes to stdout
type StdOutWriter struct{}

func (StdOutWriter) Write(msg streams.Event) {
	log.WithFields(logrus.Fields{
		"name":  "StdOutWriter",
		"event": msg,
	}).Info("writing event")
}
