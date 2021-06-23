package api

import (
	"fmt"
	"net/http"
	"time"

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

		// should closehandler send sessionEnd?
		// depends if the client can send end before dropping (like on browsers .onbeforeunload event)

		// TODO make this not hacky... use mutex lock. ensure one session END can be written
		sessionEndSent := false
		c.SetCloseHandler(func(code int, text string) error {
			log.Infof("Invoking CloseHandler code: %s, text: %s", code, text)
			if !sessionEndSent {
				endEvent := SessionEvent{
					Timestamp: int(time.Now().UnixNano() / 1000000),
					Type:      EndSession,
					SessionID: sessionID,
				}
				err = wr.Write(endEvent, sessionID)
				if err != nil {
					log.Error("session end write:", err)
					return err
				}
			}
			return nil
		})

		for {
			var message []SessionEvent
			err = c.ReadJSON(&message)

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

			// TODO this should handel serial nature of the messages
			// but it would be better to explicitly control that ...create event writer bus that has message channel per session
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
					sessionEndSent = true
					err = wr.Write(evt, sessionID)
					if err != nil {
						log.Error("write:", err)
					}
					log.Info("session ended: closing connection")

					// not sure about the etiquette here
					msg := struct {
						Message string `json:"message"`
					}{
						"bye",
					}
					c.WriteJSON(msg)
					// c.Close() is deferred
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
