package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

// SetupListenPostgreSQL setups a event listener for all events in a PostgreSQL database
// Note: this requires setting up a LISTEN/NOTIFY proc in the database before events are receivable
func SetupListenPostgreSQL(messages chan []byte, eventName string) {
	var conninfo string = fmt.Sprintf(
		"dbname=%s user=%s password=%s sslmode=disable",
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)

	_, err := sql.Open("postgres", conninfo)
	if err != nil {
		log.Fatal(err)
	}

	reportProblem := func(ev pq.ListenerEventType, err error) {
		if err != nil {
			log.Fatal("Failed to connect postgresql listener:  ", err.Error())
		}
	}

	listener := pq.NewListener(conninfo, 10*time.Second, time.Minute, reportProblem)
	err = listener.Listen(eventName)
	if err != nil {
		log.Fatal(err)
	}

	// TODO done channel to break listener
	go func() {
		for {
			select {
			case n := <-listener.Notify:
				var prettyJSON bytes.Buffer

				err := json.Indent(&prettyJSON, []byte(n.Extra), "", "\t")
				if err != nil {
					log.Fatalln("Error processing JSON: ", err)
					return
				}

				messages <- prettyJSON.Bytes()
			case <-time.After(120 * time.Second):
				log.WithFields(logrus.Fields{
					"name": "PostgresMonitor",
				}).Debug("Checking postgresql connection")

				go func() {
					if err := listener.Ping(); err != nil {
						log.Fatal(err)
					}
				}()
			}
		}
	}()
}
