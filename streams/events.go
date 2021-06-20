package streams

import "net/http"

// TODO types for start/stop to be able to do type switches
// TODO enum of known event names -- needed if they are to be written to a normalized data store

// Event is an incoming streamed event
type Event struct {
	Time      int    `json:"time"`
	Type      string `json:"type"`
	SessionID string `json:"session_id"`
	Name      string `json:"name"`
}

// EventWriter is an interface that implements writing Events to a destination
type EventWriter interface {
	Write(msg Event)
}

// EventService is a service to take incoming http events and write them to some output dest
type EventService struct {
	EventWriter
	HTTPHandler http.HandlerFunc
}
