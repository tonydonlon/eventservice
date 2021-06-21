package api

const StartSession = "SESSION_START"
const EndSession = "SESSION_END"

const TypeSession = "SESSION"
const TypeEvent = "EVENT"

// TODO types for start/stop to be able to do type switches instead of sharing this type
// Event is an incoming streamed event
type SessionEvent struct {
	Timestamp int    `json:"timestamp"`
	Type      string `json:"type"`
	SessionID string `json:"session_id"`
	Name      string `json:"name"`
}

// ClientError is an error message sent to client when request is incorrect
type ClientError struct {
	Error string `json:"error"`
}

// Event is a streamed event
type Event struct {
	Type      string `json:"type"`
	Timestamp int64  `json:"timestamp"`
	Name      string `json:"name"`
}

// SessionEventsResponse request response for a Sessions's events
type SessionEventsResponse struct {
	Type     string  `json:"type"`
	Start    int64   `json:"start"`
	End      int64   `json:"end"`
	Children []Event `json:"children"`
}
