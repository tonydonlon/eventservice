package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

func init() {
	// TODO move logger setup to a logger module
	checkLevel := func(level string) bool {
		for _, a := range log.AllLevels {
			if a.String() == level {
				return true
			}
		}
		return false
	}
	var logLevel = os.Getenv("LOG_LEVEL")
	if !checkLevel(logLevel) {
		log.Infof("Unknown LOG_LEVEL: '%s' using default 'info'", logLevel)
		logLevel = "info"
	}
	log.SetFormatter(&log.JSONFormatter{
		PrettyPrint: logLevel == "debug",
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

var upgrader = websocket.Upgrader{}

func event(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			// TODO distinguish disconnect from real error
			log.Info(err)
			break
		}
		log.Debugf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Error("write:", err)
			break
		}
	}
}

func main() {
	log.Info("eventservice")
	http.HandleFunc("/event", event)
	var portNumber = os.Getenv("WS_PORT")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%s", portNumber), nil))
}
