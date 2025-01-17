package src

import (
	"database/sql"
	"github.com/gorilla/mux"
	"quit4real.today/src/api"
	"quit4real.today/src/cron"
	"quit4real.today/src/endpoint"
	"quit4real.today/src/handler/command"
	"quit4real.today/src/handler/query"
	"quit4real.today/src/repository"
)

type App struct {
	DatabaseImpl    *repository.DatabaseImpl
	Endpoints       *endpoint.Endpoints
	CommandHandlers *CommandHandlers
	QueryHandlers   *QueryHandlers
	Repositories    *Repositories
	Jobs            *cron.Jobs
	SteamApi        *api.SteamApi
}

type CommandHandlers struct {
	FailsCommandHandler   *command.FailsCommandHandler
	TrackerCommandHandler *command.TrackerCommandHandler
	UserCommandHandler    *command.UserCommandHandler
}

type QueryHandlers struct {
	FailsQueryHandler *query.FailQueryHandler
	UserQueryHandler  *query.UserQueryHandler
}

type Repositories struct {
	FailRepository    *repository.FailRepository
	UserRepository    *repository.UserRepository
	TrackerRepository *repository.TrackerRepository
}

func AppInit(dataBaseConnection *sql.DB) *App {
	var app = App{}
	databaseImpl := repository.DatabaseImpl{DB: dataBaseConnection}
	repositories := createRepositories(app)
	commandHandler := createCommandHandlers(app)
	queryHandler := createQueryHandlers(app)
	jobs := createJobs(app)
	steamApi := createSteamApi()
	endpoints := createEndpoints(app)

	app.DatabaseImpl = &databaseImpl
	app.Repositories = &repositories
	app.CommandHandlers = &commandHandler
	app.QueryHandlers = &queryHandler
	app.Jobs = &jobs
	app.SteamApi = &steamApi
	app.Endpoints = &endpoints
	return &app
}

func createRepositories(app App) Repositories {
	trackerRepo := repository.TrackerRepository{DatabaseImpl: app.DatabaseImpl}
	userRepo := repository.UserRepository{DatabaseImp: app.DatabaseImpl}
	failRepo := repository.FailRepository{DatabaseImpl: app.DatabaseImpl}

	return Repositories{
		TrackerRepository: &trackerRepo,
		UserRepository:    &userRepo,
		FailRepository:    &failRepo,
	}
}

func createCommandHandlers(app App) CommandHandlers {
	failCommandHandler := command.FailsCommandHandler{FailRepository: app.Repositories.FailRepository}
	trackerCommandHandler := command.TrackerCommandHandler{TrackerRepository: app.Repositories.TrackerRepository}
	userCommandHandler := command.UserCommandHandler{UserRepository: app.Repositories.UserRepository}

	return CommandHandlers{
		FailsCommandHandler:   &failCommandHandler,
		TrackerCommandHandler: &trackerCommandHandler,
		UserCommandHandler:    &userCommandHandler,
	}
}

func createQueryHandlers(app App) QueryHandlers {
	failQueryHandler := query.FailQueryHandler{FailRepository: app.Repositories.FailRepository}
	userQueryHandler := query.UserQueryHandler{UserRepository: app.Repositories.UserRepository}

	return QueryHandlers{
		FailsQueryHandler: &failQueryHandler,
		UserQueryHandler:  &userQueryHandler,
	}
}

func createJobs(app App) cron.Jobs {
	failCron := cron.FailCron{
		UserQueryHandler:      app.QueryHandlers.UserQueryHandler,
		TrackerCommandHandler: app.CommandHandlers.TrackerCommandHandler,
	}

	return cron.Jobs{
		FailCron: &failCron,
	}
}

func createSteamApi() api.SteamApi {
	return api.SteamApi{}
}

func createEndpoints(app App) endpoint.Endpoints {
	router := mux.NewRouter()

	userEndpoint := endpoint.UserEndpoint{
		Router:                router,
		UserCommandHandler:    app.CommandHandlers.UserCommandHandler,
		TrackerCommandHandler: app.CommandHandlers.TrackerCommandHandler,
	}

	failEndpoint := endpoint.FailEndpoint{
		Router:           router,
		FailQueryHandler: app.QueryHandlers.FailsQueryHandler,
	}

	return endpoint.Endpoints{
		Router:       router,
		UserEndpoint: &userEndpoint,
		FailEndpoint: &failEndpoint,
	}
}
