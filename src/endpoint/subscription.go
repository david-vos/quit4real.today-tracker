package endpoint

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"quit4real.today/logger"
	"quit4real.today/src/handler/command"
	"quit4real.today/src/model"
)

type SubscriptionEndpoint struct {
	Router                     *mux.Router
	SubscriptionCommandHandler *command.SubscriptionCommandHandler
}

// Subscription handles the subscription-related endpoints.
func (endpoint *SubscriptionEndpoint) Subscription() {
	logger.Info("Starting subscription endpoints")
	endpoint.Router.HandleFunc("/subscriptions", endpoint.AddSubscription()).Methods("POST")
	logger.Info("Subscription endpoints started")
}

// AddSubscription handles adding a subscription for a user.
func (endpoint *SubscriptionEndpoint) AddSubscription() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Got request to add a new Subscription")
		if r.Method != http.MethodPost {
			logger.Debug("Not a POST request")
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		// Parse JSON request body
		var subscription model.Subscription
		err := json.NewDecoder(r.Body).Decode(&subscription)
		if err != nil {
			logger.Debug("Error decoding subscription: " + err.Error())
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		err = endpoint.SubscriptionCommandHandler.Add(subscription)
		if err != nil {
			logger.Debug("Error adding subscription: " + err.Error())
			http.Error(w, "Error adding subscription: "+err.Error(), http.StatusInternalServerError)
		}

		// Respond with success
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte("Subscription created successfully"))
		if err != nil {
			return
		}
	}
}
