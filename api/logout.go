package api

import (
	"github.com/sirupsen/logrus"
)

// StdOutWriter implements events.EventWriter and writes to stdout
type StdOutWriter struct{}

func (s *StdOutWriter) Write(msg SessionEvent, sessionID string) error {
	log.WithFields(logrus.Fields{
		"name":  "StdOutWriter",
		"event": msg,
	}).Info("writing event")
	return nil
}

func (s *StdOutWriter) Init() error {
	log.WithFields(logrus.Fields{
		"name": "StdOutWriter",
	}).Info("Init")
	return nil
}

// StdOutReader implements SessionReader interface
type StdOutReader struct{}

func (StdOutReader) SessionEvents(sessionID string) (*SessionEventsResponse, error) {
	return nil, nil
}

func (StdOutReader) Init() error {
	log.WithFields(logrus.Fields{
		"name": "StdOutWriter",
	}).Info("Init")
	return nil
}
