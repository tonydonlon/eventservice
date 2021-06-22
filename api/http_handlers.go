package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// GetSessionEvents gets a session and all it's events over REST
func GetSessionEvents(reader SessionReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		// TODO is this the best way to validate UUID?
		sessionID, error := uuid.FromString(vars["session_id"])
		if error != nil {
			log.Errorf("%s is not UUID", vars["session_id"])
			http.Error(w, "session_id must be a UUID", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		log.Infof("Retrieving events for sessionID: %s", sessionID)
		sessionEvents, err := reader.SessionEvents(sessionID.String())
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError)
			return
		}

		// no results for sessionID
		if sessionEvents == nil {
			log.Infof("no events fountd for %s", sessionID)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(sessionEvents)
	}

}
