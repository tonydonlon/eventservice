package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/tonydonlon/eventservice/api"
	"github.com/tonydonlon/eventservice/logger"
)

// TODO remove global
var log *logrus.Logger

func init() {
	log = logger.GetLogger()
}

// TODO api.EventService should be composed of all the handlers
func SessionEventHandler(reader api.SessionReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		// TODO is this the best way to validate UUID?
		sessionID, error := uuid.FromString(vars["session_id"])
		if error != nil {
			log.Errorf("%s is not UUID", vars["session_id"])
			http.Error(w, "session_id must be a UUID", http.StatusBadRequest)
			return
		}
		log.Infof("Retrieving events for sessionID: %s", sessionID)
		w.Header().Set("Content-Type", "application/json")

		sessionEvents, err := reader.SessionEvents(sessionID.String())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}

		// no results for sessionID
		if sessionEvents == nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(sessionEvents)
	}

}
