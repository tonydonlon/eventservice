package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/tonydonlon/eventservice/api"
	"github.com/tonydonlon/eventservice/logger"
	"github.com/tonydonlon/eventservice/websocket"
	"github.com/tonydonlon/eventservice/writers"
)

var log *logrus.Logger

func init() {
	log = logger.GetLogger()
	// TODO validate config
}

func main() {
	ws := websocket.WebsocketEventHandler{}
	srv := &api.EventService{
		HTTPHandler: ws.Handler(),
		EventWriter: writers.StdOutWriter{},
	}

	http.HandleFunc("/event", srv.HTTPHandler)
	// TODO REST/stateless http endpoint for fallback if WS is not supported client-side
	var portNumber = os.Getenv("HTTP_PORT")
	log.Info("eventservice", fmt.Sprintf(" listening on localhost:%s", portNumber))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%s", portNumber), nil))
}
