package endpoint

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"quit4real.today/logger"
	"quit4real.today/src/handler/query"
)

type GamesEndpoint struct {
	Router           *mux.Router
	GameQueryHandler *query.GameQueryHandler
}

func (endpoint *GamesEndpoint) Games() {
	logger.Info("Starting Games Endpoint")
	endpoint.Router.HandleFunc("/games/{searchParam}/{platformId}", endpoint.SearchGames()).Methods("GET")
	logger.Info("Games Endpoint Started")
}

func (endpoint *GamesEndpoint) SearchGames() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		searchParam := vars["searchParam"]
		platformId := vars["platformId"]
		if platformId != "steam" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		games, err := endpoint.GameQueryHandler.Search(searchParam, "steam")
		if err != nil {
			logger.Fail(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			_, err := w.Write([]byte(err.Error()))
			if err != nil {
				return
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(games); err != nil {
			logger.Fail(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
