package main

import (
	"project/config"
	"project/logger"
	"project/main/db"
	"project/main/endpoint"
)

func main() {
	// Setup logger for error handling
	logger.Info("Starting the logger")
	config.InitLogger()
	logger.Info("Successfully initialized the logger")

	// Setup sql lite database
	logger.Info("Starting the database")
	dataBaseConnection := db.Setup()
	defer func() {
		err := dataBaseConnection.Close()
		if err != nil {
			logger.Fail("Failed to close database" + err.Error())
		}
	}()
	logger.Info("Successfully started DB connection")

	// TODO: setup a way to init the dependency injections
	var app *App = AppInit(dataBaseConnection)

	// Setup http server
	logger.Info("Starting http server")
	endpoint.Init(app.Endpoints)
	logger.Info("Successfully initialized the http server")

	select {}
}
