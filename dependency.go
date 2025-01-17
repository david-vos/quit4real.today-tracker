package main

import (
	"database/sql"
	"project/handlers"
	"project/main/api"
	"project/repository"
)

// DependencyInit initializes the UserController with the necessary dependencies.
func DependencyInit(dbc *sql.DB) (*handlers.UserController, *api.CronController, *handlers.FailsController) {
	//Database Layer Controller definition
	databaseControllerImp := &repository.DatabaseController{
		DB: dbc,
	}
	// Repo Layer Controller definition
	userRepoControllerImp := &repository.UserRepoController{
		DbContr: databaseControllerImp,
	}
	trackerRepoControllerImp := &repository.TrackerRepoController{
		DbContr: databaseControllerImp,
	}
	failedRepoControllerImp := &repository.FailedRepoController{
		DbContr: databaseControllerImp,
	}

	// Top layer Controller Definition
	userHandlerControllerImp := &handlers.UserController{
		UserRepoContr:    userRepoControllerImp,
		TrackerRepoContr: trackerRepoControllerImp,
	}
	failsController := &handlers.FailsController{
		FailRepoContr: failedRepoControllerImp,
		UserRepoContr: userRepoControllerImp,
	}
	cronControllerImp := &api.CronController{
		UserRepoContr:    userRepoControllerImp,
		TrackerRepoContr: trackerRepoControllerImp,
		FailsContr:       failsController,
	}

	return userHandlerControllerImp, cronControllerImp, failsController
}
