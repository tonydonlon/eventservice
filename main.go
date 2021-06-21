package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/tonydonlon/eventservice/api"
	"github.com/tonydonlon/eventservice/handlers"
	"github.com/tonydonlon/eventservice/logger"
	"github.com/tonydonlon/eventservice/writers"
)

var log *logrus.Logger

func init() {
	log = logger.GetLogger()
	// TODO validate config
}

func main() {
	var eventWriter api.EventWriter

	// TODO move all this to api.EventService factory
	if os.Getenv("DATASTORE") == "postgres" {
		eventWriter = &writers.PostgreSQLWriter{}
	} else {
		eventWriter = &writers.StdOutWriter{}
	}

	if err := eventWriter.Init(); err != nil {
		log.Error(err)
		return
	}

	service := &api.EventService{
		WebsocketHandler: handlers.WebsocketHandler(eventWriter),
		EventWriter:      eventWriter,
	}
	// TODO REST/stateless http endpoint for fallback if WS is not supported client-side
	var portNumber = os.Getenv("HTTP_PORT")
	serviceAddress := fmt.Sprintf("localhost:%s", portNumber)

	router := mux.NewRouter()
	router.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("healthcheck")
		io.WriteString(w, "OK")
	})
	router.HandleFunc("/ws", service.WebsocketHandler)
	reader := &writers.PostgreSQLReader{}
	reader.Init()
	router.HandleFunc("/session/{session_id}", handlers.SessionEventHandler(reader)).Methods("GET")

	srv := &http.Server{
		Handler:      router,
		Addr:         serviceAddress,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// TODO graceful shutdown
	log.Info("eventservice", fmt.Sprintf(" listening on %s", serviceAddress))
	log.Fatal(srv.ListenAndServe())

	// TODO gRPC server on different port
}
