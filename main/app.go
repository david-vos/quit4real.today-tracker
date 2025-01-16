package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"project/main/endpoint"
)

type App struct {
	DB        *sql.DB
	Endpoints *endpoint.Endpoints
}

func AppInit(dataBaseConnection *sql.DB) *App {
	return &App{
		DB:        dataBaseConnection,
		Endpoints: createEndpoints(),
	}

}

func createEndpoints() *endpoint.Endpoints {
	router := mux.NewRouter()

	userEndpoint := endpoint.UserEndpoint{
		Router: router,
	}

	failEndpoint := endpoint.FailEndpoint{
		Router: router,
	}

	endpoints := endpoint.Endpoints{
		Router:       router,
		UserEndpoint: &userEndpoint,
		FailEndpoint: &failEndpoint,
	}
	return &endpoints
}
