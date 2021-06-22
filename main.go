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
	"github.com/tonydonlon/eventservice/broker"
	"github.com/tonydonlon/eventservice/logger"
)

var log *logrus.Logger

func init() {
	log = logger.GetLogger()
	// TODO validate config
}

func main() {
	service, err := api.NewEventService()
	if err != nil {
		log.Fatal(err)
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
	router.HandleFunc("/session/{session_id}", service.SessionEventHandler).Methods("GET")

	// for monitoring the writes to the database; TODO endpoint to stream every notification
	databaseBus := broker.NewEventBus("Database")
	done := make(chan bool)
	databaseBus.Start(done)
	api.SetupListenPostgreSQL(databaseBus.Messages, "customer_events")

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
