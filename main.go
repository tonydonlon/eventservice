package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/tonydonlon/eventservice/logger"
	"github.com/tonydonlon/eventservice/websocket"
)

var log *logrus.Logger

func init() {
	log = logger.GetLogger()
}

func main() {
	log.Info("eventservice")
	http.HandleFunc("/event", websocket.EventHandler)
	var portNumber = os.Getenv("WS_PORT")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%s", portNumber), nil))
}
