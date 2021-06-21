package api

import "net/http"

// TODO types for start/stop to be able to do type switches

const StartSession = "SESSION_START"
const EndSession = "SESSION_END"

// Event is an incoming streamed event
type Event struct {
	Time      int    `json:"time"`
	Type      string `json:"type"`
	SessionID string `json:"session_id"`
	Name      string `json:"name"`
}

// TODO Factory NewEventService

// TODO enum of known event names -- needed if they are to be written to a normalized data store
// EventNames POLL db for new event names; mutex lock serice memory mapped lookup
type EventNames struct {
}

// EventWriter is an interface that implements writing Events to a destination
type EventWriter interface {
	Write(msg Event, sessionID string) error
}

// EventService is a service to take incoming http events and write them to some output dest
type EventService struct {
	EventWriter
	HTTPHandler      http.HandlerFunc
	WebsocketHandler http.HandlerFunc
}

// ClientError is an error message sent to client when request is incorrect
type ClientError struct {
	Error string `json:"error"`
}
