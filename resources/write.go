package resources

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"

	"github.com/Financial-Times/list-notifications-rw/db"
	"github.com/Financial-Times/list-notifications-rw/mapping"
	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

// WriteNotification will write a new notification for the provided list.
func WriteNotification(dumpRequests bool, mapper mapping.NotificationsMapper, db db.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if dumpRequests {
			dumpRequest(r)
		}

		decoder := json.NewDecoder(r.Body)
		uuid := mux.Vars(r)["uuid"]

		notification, err := mapper.MapRequestToInternalNotification(uuid, decoder)
		if err != nil {
			logrus.WithError(err).
				WithField("uuid", uuid).
				WithField("tid", r.Header.Get("X-Request-Id")).
				Error("Invalid request! See error for details.")
			writeMessage("Invalid Request body.", 400, w)
			return
		}

		tx, err := db.Open()
		if err != nil {
			logrus.WithError(err).Error("Failed to connect to mongo")
			writeMessage("An internal server error prevented processing of your request.", 500, w)
			return
		}

		defer tx.Close()

		tx.WriteNotification(notification)

		logrus.WithField("uuid", uuid).WithField("transaction_id", r.Header.Get("X-Request-Id")).Info("Successfully processed a notification for this list.")
		w.WriteHeader(200)
	}
}

func dumpRequest(r *http.Request) {
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		logrus.WithError(err).Warn("Failed to dump request!")
		return
	}
	logrus.Info(string(dump))
}

type msg struct {
	Message string `json:"message"`
}

func writeMessage(message string, status int, w http.ResponseWriter) {
	w.WriteHeader(status)

	m := msg{Message: message}
	encoder := json.NewEncoder(w)
	encoder.Encode(m)

	w.Header().Add("Content-Type", "application/json")
}
