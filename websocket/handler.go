package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/tonydonlon/eventservice/logger"
	"github.com/tonydonlon/eventservice/streams"
	"github.com/tonydonlon/eventservice/writers"
)

var log *logrus.Logger

func init() {
	log = logger.GetLogger()
}

type WebsocketEventHandler struct{}

// Handler handles incoming websocket events
func (ws *WebsocketEventHandler) Handler() http.HandlerFunc {
	// TODO handle sessionID as a parameter
	return func(w http.ResponseWriter, r *http.Request) {
		var upgrader = websocket.Upgrader{}
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("upgrade:", err)
			return
		}
		defer c.Close()
		wr := writers.StdOutWriter{}
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

			// TODO need a type switch here
			for _, evt := range message {
				go wr.Write(evt)
				if evt.Type == "SESSION_END" {
					log.Info("session ended")
					c.Close()
					return
				}
			}

			if err != nil {
				log.Error("write:", err)
				break
			}
		}
	}
}
