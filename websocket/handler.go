package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/tonydonlon/eventservice/logger"
	"github.com/tonydonlon/eventservice/streams"
)

var log *logrus.Logger

func init() {
	log = logger.GetLogger()
}

var upgrader = websocket.Upgrader{}

// EventHandler handles incoming websocket events
func EventHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		var message []streams.Event
		err = c.ReadJSON(&message)
		// TODO message validation;missing values, etc.
		if err != nil {
			log.Error("Error reading json.", err)
		}

		if err != nil {
			log.Error(err)
			break
		}

		// SESSION_START is guarateed to be first and SESSION_END is guaranteed to be last
		// TODO how to be defensive if that is not the case
		// spew.Dump(len(message), message[0], message[len(message)-1])

		log.Debug(message)
		for _, evt := range message {
			if evt.Type == "SESSION_END" {
				log.Info("SESSION_END")
				c.Close()
				return
			}
		}

		err = c.WriteJSON((message))
		if err != nil {
			log.Error("write:", err)
			break
		}
	}
}
