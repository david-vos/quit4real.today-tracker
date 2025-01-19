package endpoint

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"quit4real.today/logger"
	"quit4real.today/src/handler/query"
)

type FailEndpoint struct {
	Router           *mux.Router
	FailQueryHandler *query.FailQueryHandler
}

func (endpoint *FailEndpoint) Fail() {
	logger.Info("Trying to start fail endpoint")
	endpoint.Router.HandleFunc("/fail/leaderboard", endpoint.getLeaderboard()).Methods("GET")
	logger.Info("Fail endpoint started")
}

func (endpoint *FailEndpoint) getLeaderboard() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		logger.Info("Getting leaderboard for this endpoint")
		failsLeaderBoard, err := endpoint.FailQueryHandler.GetLeaderBoard()
		if err != nil {
			logger.Fail(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return // Ensure to return after writing the header
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(failsLeaderBoard); err != nil {
			logger.Fail(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
