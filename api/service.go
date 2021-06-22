package api

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/tonydonlon/eventservice/broker"
	"github.com/tonydonlon/eventservice/logger"
)

var log *logrus.Logger

func init() {
	log = logger.GetLogger()
}

// Event types are all the allowable/known event types from the normalized data store
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
// TODO list of listeners for all configured/supported protocols mapped to routes

// EventService is a service to take incoming http events and write them to some output dest
type EventService struct {
	SessionEventHandler http.HandlerFunc
	WebsocketHandler    http.HandlerFunc
	EventTypes          EventTypes
	EventBus            broker.EventBus
	DatabaseBus         broker.EventBus
	SSEHandler          http.HandlerFunc
}

// NewEventService creates the EventService API
func NewEventService() (*EventService, error) {

	// TODO make these reader/writer abstractions go away
	var eventWriter EventWriter
	if os.Getenv("DATASTORE") == "postgres" {
		eventWriter = &PostgreSQLWriter{}
	} else {
		eventWriter = &StdOutWriter{}
	}
	if err := eventWriter.Init(); err != nil {
		log.Error(err)
		return nil, err
	}

	var eventReader SessionReader
	if os.Getenv("DATASTORE") == "postgres" {
		eventReader = &PostgreSQLReader{}
	} else {
		eventReader = &StdOutReader{}
	}
	if err := eventReader.Init(); err != nil {
		log.Error(err)
		return nil, err
	}

	// TODO EventTypes could refresh on a POLL to db for new event names; mutex lock serice memory mapped struct
	// SELECT event_name, event_name_id FROM event_names
	eventTypes := EventTypes{
		"test":   1,
		"event1": 2,
		"event2": 3,
	}

	// for monitoring the writes to the database
	databaseBus := broker.NewEventBus("Database")
	done := make(chan bool)
	databaseBus.Start(done)
	SetupListenPostgreSQL(databaseBus.Messages, "customer_events")

	service := &EventService{
		WebsocketHandler:    WebsocketHandler(eventWriter),
		SessionEventHandler: GetSessionEvents(eventReader),
		SSEHandler:          SSEHandler(databaseBus),
		EventBus:            *broker.NewEventBus("EventService"),
		DatabaseBus:         *databaseBus,
		EventTypes:          eventTypes,
	}

	return service, nil
}
