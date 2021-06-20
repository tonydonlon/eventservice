package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/tonydonlon/eventservice/logger"
	"github.com/tonydonlon/eventservice/streams"
	"github.com/tonydonlon/eventservice/websocket"
	"github.com/tonydonlon/eventservice/writers"
)

var log *logrus.Logger

func init() {
	log = logger.GetLogger()
}

func main() {
	log.Info("eventservice")

	ws := websocket.WebsocketEventHandler{}
	srv := &streams.EventService{
		HTTPHandler: ws.Handler(),
		EventWriter: writers.StdOutWriter{},
	}

	http.HandleFunc("/event", srv.HTTPHandler)
	var portNumber = os.Getenv("WS_PORT")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%s", portNumber), nil))
}
