package main

import (
	"database/sql"
	"project/handlers"
	"project/repository"
)

// DependencyInit initializes the UserController with the necessary dependencies.
func DependencyInit(dbc *sql.DB) (*handlers.UserController, *handlers.CronController) {
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

	// Top layer Controller Definition
	userHandlerControllerImp := &handlers.UserController{
		UserRepoContr:    userRepoControllerImp,
		TrackerRepoContr: trackerRepoControllerImp,
	}
	cronControllerImp := &handlers.CronController{
		UserRepoContr:    userRepoControllerImp,
		TrackerRepoContr: trackerRepoControllerImp,
	}

	return userHandlerControllerImp, cronControllerImp
}
