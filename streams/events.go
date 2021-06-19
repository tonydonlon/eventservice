package streams

// TODO types for start/stop to be able to do type switches
// TODO enum of known event names -- needed if they are to be written to a normalized data store

// Event is an incoming streamed event
type Event struct {
	Time      int    `json:"time"`
	Type      string `json:"type"`
	SessionID string `json:"session_id"`
	Name      string `json:"name"`
}
