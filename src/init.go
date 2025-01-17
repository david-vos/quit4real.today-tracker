package src

import (
	"quit4real.today/config"
	"quit4real.today/logger"
	"quit4real.today/src/db"
	"quit4real.today/src/endpoint"
)

func StartApp() {
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

	logger.Info("Starting the dependency init")
	var app *App = AppInit(dataBaseConnection)
	logger.Info("Successfully started dependency init")

	// Setup http server
	logger.Info("Starting http server")
	endpoint.Init(app.Endpoints)
	logger.Info("Successfully initialized the http server")

	logger.Info("Starting the cronJobs")
	app.Jobs.StartAll()
	logger.Info("Successfully started cronJobs")

	select {}
}
