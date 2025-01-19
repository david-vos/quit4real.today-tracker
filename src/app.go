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

// App represents the main application structure.
type App struct {
	DatabaseImpl    *repository.DatabaseImpl
	Endpoints       *endpoint.Endpoints
	CommandHandlers *CommandHandlers
	QueryHandlers   *QueryHandlers
	Repositories    *Repositories
	Jobs            *cron.Jobs
	SteamApi        *api.SteamApi
}

// CommandHandlers holds all command handlers.
type CommandHandlers struct {
	FailsCommandHandler   *command.FailsCommandHandler
	TrackerCommandHandler *command.TrackerCommandHandler
	UserCommandHandler    *command.UserCommandHandler
}

// QueryHandlers holds all query handlers.
type QueryHandlers struct {
	FailsQueryHandler *query.FailQueryHandler
	UserQueryHandler  *query.UserQueryHandler
}

// Repositories holds all repositories.
type Repositories struct {
	FailRepository    *repository.FailRepository
	UserRepository    *repository.UserRepository
	TrackerRepository *repository.TrackerRepository
}

// AppInit initializes the application with the provided database connection.
func AppInit(dataBaseConnection *sql.DB) *App {
	databaseImpl := &repository.DatabaseImpl{DB: dataBaseConnection}
	repositories := createRepositories(databaseImpl)
	steamApi := createSteamApi()
	commandHandlers := createCommandHandlers(repositories, steamApi)
	queryHandlers := createQueryHandlers(repositories)
	jobs := createJobs(queryHandlers, commandHandlers)
	endpoints := createEndpoints(commandHandlers, queryHandlers, steamApi)

	return &App{
		DatabaseImpl:    databaseImpl,
		Repositories:    repositories,
		CommandHandlers: commandHandlers,
		QueryHandlers:   queryHandlers,
		Jobs:            jobs,
		SteamApi:        steamApi,
		Endpoints:       endpoints,
	}
}

func createRepositories(databaseImpl *repository.DatabaseImpl) *Repositories {
	return &Repositories{
		TrackerRepository: &repository.TrackerRepository{DatabaseImpl: databaseImpl},
		UserRepository:    &repository.UserRepository{DatabaseImpl: databaseImpl},
		FailRepository:    &repository.FailRepository{DatabaseImpl: databaseImpl},
	}
}

func createCommandHandlers(repositories *Repositories, steamApi *api.SteamApi) *CommandHandlers {
	failsHandler := &command.FailsCommandHandler{FailRepository: repositories.FailRepository}
	trackerCommandHandler := &command.TrackerCommandHandler{
		SteamApi:            steamApi,
		TrackerRepository:   repositories.TrackerRepository,
		FailsCommandHandler: failsHandler,
	}

	return &CommandHandlers{
		FailsCommandHandler:   failsHandler,
		TrackerCommandHandler: trackerCommandHandler,
		UserCommandHandler:    &command.UserCommandHandler{UserRepository: repositories.UserRepository},
	}
}

func createQueryHandlers(repositories *Repositories) *QueryHandlers {
	return &QueryHandlers{
		FailsQueryHandler: &query.FailQueryHandler{FailRepository: repositories.FailRepository},
		UserQueryHandler:  &query.UserQueryHandler{UserRepository: repositories.UserRepository},
	}
}

func createJobs(queryHandlers *QueryHandlers, commandHandlers *CommandHandlers) *cron.Jobs {

	return &cron.Jobs{
		FailCron: &cron.FailCron{
			UserQueryHandler:      queryHandlers.UserQueryHandler,
			TrackerCommandHandler: commandHandlers.TrackerCommandHandler,
		},
	}
}

func createSteamApi() *api.SteamApi {
	return &api.SteamApi{}
}

func createEndpoints(commandHandlers *CommandHandlers, queryHandlers *QueryHandlers, steamApi *api.SteamApi) *endpoint.Endpoints {
	router := mux.NewRouter()

	return &endpoint.Endpoints{
		Router: router,
		UserEndpoint: &endpoint.UserEndpoint{
			Router:                router,
			SteamApi:              steamApi,
			UserCommandHandler:    commandHandlers.UserCommandHandler,
			TrackerCommandHandler: commandHandlers.TrackerCommandHandler,
		},
		FailEndpoint: &endpoint.FailEndpoint{
			Router:           router,
			FailQueryHandler: queryHandlers.FailsQueryHandler,
		},
	}
}
