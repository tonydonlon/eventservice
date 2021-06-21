package api

import "net/http"

// TODO map of known event names -- needed if they are to be written to a normalized data store
// SELECT event_name, event_name_id FROM event_names
// EventTypes could refresh on a POLL to db for new event names; mutex lock serice memory mapped struct
type EventTypes map[string]int

// EventWriter is an interface that implements writing Events to a destination
type EventWriter interface {
	Write(msg SessionEvent, sessionID string) error
	Init() error
}

// SessionReader retrieves a session and all related events
type SessionReader interface {
	SessionEvents(sessionID string) (*SessionEventsResponse, error)
	Init() error
}

// TODO logrus.Logger service singleton logger?

// EventService is a service to take incoming http events and write them to some output dest
type EventService struct {
	EventWriter
	SessionReader
	// TODO should be a list of listeners for all configured/supported protocols
	HTTPHandler      http.HandlerFunc
	WebsocketHandler http.HandlerFunc
	EventTypes       EventTypes
}

// TODO Factory NewEventService
