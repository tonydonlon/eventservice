package broker

import (
	"github.com/sirupsen/logrus"
	"github.com/tonydonlon/eventservice/logger"
)

var log *logrus.Logger

func init() {
	log = logger.GetLogger()
}

// EventBus is a simple pub sub event bus
type EventBus struct {
	Messages      chan []byte
	Subscribe     chan chan []byte
	Unsubscribe   chan chan []byte
	subscriptions map[chan []byte]bool
	halt          chan<- bool
	Name          string
}

// NewEventBus creates a new EventBus
func NewEventBus(name string) *EventBus {
	bus := &EventBus{
		Messages:      make(chan []byte, 1),
		Subscribe:     make(chan chan []byte),
		Unsubscribe:   make(chan chan []byte),
		subscriptions: make(map[chan []byte]bool),
		Name:          name,
	}
	return bus
}

// Start makes the EventBus start listening
func (b *EventBus) Start(done chan<- bool) {
	b.halt = done
	log.Infof("Starting %s EventBus", b.Name)
	go b.listen()
}

// ConnectionCount returns the number of active client connections
func (b *EventBus) ConnectionCount() int {
	return len(b.subscriptions)
}

func (b *EventBus) listen() {
	for {
		select {
		case chnl := <-b.Subscribe:
			// TODO maybe just use uuid for map key?
			b.subscriptions[chnl] = true
			log.WithFields(logrus.Fields{
				"name":    "EventBus",
				"busName": b.Name,
			}).Infof("Eventbus %s: Client added: %v. %d registered clients ", b.Name, chnl, b.ConnectionCount())
		case chnl := <-b.Unsubscribe:
			delete(b.subscriptions, chnl)
			log.WithFields(logrus.Fields{
				"name":    "EventBus",
				"busName": b.Name,
			}).Infof("Eventbus %s: Removed client: %v. %d registered clients", b.Name, chnl, b.ConnectionCount())
		case event := <-b.Messages:
			for subscriberChannel := range b.subscriptions {
				subscriberChannel <- event
			}
		}

	}
}
