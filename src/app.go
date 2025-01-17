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
	databaseImpl := &repository.DatabaseImpl{DB: dataBaseConnection}
	repositories := createRepositories(databaseImpl)
	commandHandlers := createCommandHandlers(repositories)
	queryHandlers := createQueryHandlers(repositories)
	jobs := createJobs(queryHandlers, commandHandlers)
	steamApi := createSteamApi()
	endpoints := createEndpoints(commandHandlers, queryHandlers)

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

func createCommandHandlers(repositories *Repositories) *CommandHandlers {
	return &CommandHandlers{
		FailsCommandHandler:   &command.FailsCommandHandler{FailRepository: repositories.FailRepository},
		TrackerCommandHandler: &command.TrackerCommandHandler{TrackerRepository: repositories.TrackerRepository},
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

func createEndpoints(commandHandlers *CommandHandlers, queryHandlers *QueryHandlers) *endpoint.Endpoints {
	router := mux.NewRouter()

	return &endpoint.Endpoints{
		Router: mux.NewRouter(),
		UserEndpoint: &endpoint.UserEndpoint{
			Router:                router,
			UserCommandHandler:    commandHandlers.UserCommandHandler,
			TrackerCommandHandler: commandHandlers.TrackerCommandHandler,
		},
		FailEndpoint: &endpoint.FailEndpoint{
			Router:           router,
			FailQueryHandler: queryHandlers.FailsQueryHandler,
		},
	}
}
