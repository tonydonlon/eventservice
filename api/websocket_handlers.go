package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebsocketHandler handles incoming websocket event streams
func WebsocketHandler(wr EventWriter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.URL.Query().Get("sessionId")

		var upgrader = websocket.Upgrader{}
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Error("upgrade:", err)
			return
		}
		defer c.Close()

		for {
			var message []SessionEvent
			err = c.ReadJSON(&message)

			// TODO closehandler send sessionEND

			// TODO message validation;missing values, etc.
			if err != nil {
				log.Error("Error reading json.", err)
			}

			if err != nil {
				log.Error(err)
				break
			}

			// SESSION_START is guaranteed to be first and SESSION_END is guaranteed to be last
			// TODO how to be defensive if that is not the case
			log.Warn("----------RECEIVED %d messages", len(message))
			// TODO create event writer bus that has message channel per session
			for _, evt := range message {
				switch evt.Type {
				case StartSession:
					// start should block if db schema requires sessionID existence

					// assert event message matches the one passed as parameter
					if evt.SessionID != sessionID {
						// send 400ish message and close
						msg := fmt.Sprintf("event sessionID does not match sessionID parameter: %s", sessionID)
						log.Error(msg)
						errorMsg := ClientError{msg}
						c.WriteJSON(errorMsg)
						c.Close()
						return
					}

					err = wr.Write(evt, sessionID)
					if err != nil {
						log.Error("write:", err)
					}
				case EndSession:
					err = wr.Write(evt, sessionID)
					if err != nil {
						log.Error("write:", err)
					}
					log.Info("session ended: closing connection")
					// TODO send message "bye"
					c.Close()
					return
				default:
					// TODO there needs to be an error channel for the async version of this
					//go wr.Write(evt)
					err = wr.Write(evt, sessionID)
					if err != nil {
						log.Error("write:", err)
					}
				}
			}
		}
	}
}
