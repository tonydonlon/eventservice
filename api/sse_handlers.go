package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/tonydonlon/eventservice/broker"
)

// SSEHandler returns a http handler that does Server-Sent-Events
func SSEHandler(b *broker.EventBus) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		flusher, ok := rw.(http.Flusher)
		if !ok {
			http.Error(rw, "Streaming not supported", http.StatusInternalServerError)
			return
		}

		rw.Header().Set("Content-Type", "text/event-stream")
		rw.Header().Set("Cache-Control", "no-cache")
		rw.Header().Set("Connection", "keep-alive")
		rw.Header().Set("Access-Control-Allow-Origin", "*")

		messages := make(chan []byte)

		// send channel to subscribe...think of it like passing a callback
		b.Subscribe <- messages

		// listen for when client closes; in that case, unsubscribe this handler
		notify := req.Context().Done()
		go func() {
			<-notify
			if req.Context().Err() != context.Canceled {
				log.Printf("Error with SSE handler ending: %s", req.Context().Err().Error())
			}
			b.Unsubscribe <- messages
		}()

		for {
			rw.Write(formatSSE("message", string(<-messages)))
			//fmt.Fprint(rw, <-messages)
			flusher.Flush()
		}
	}
}

func formatSSE(event string, data string) []byte {
	eventPayload := "event: " + event + "\n"
	dataLines := strings.Split(data, "\n")
	for _, line := range dataLines {
		eventPayload = eventPayload + "data: " + line + "\n"
	}
	return []byte(eventPayload + "\n")
}
