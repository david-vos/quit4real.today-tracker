package src

import (
	"database/sql"
	"github.com/gorilla/mux"
	"quit4real.today/src/cron"
	"quit4real.today/src/endpoint"
	"quit4real.today/src/handler/command"
	"quit4real.today/src/handler/query"
	"quit4real.today/src/handler/service"
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
	Services        *Services
}

type Services struct {
	SteamService *service.SteamService
	AuthService  *service.AuthService
}

// CommandHandlers holds all command handlers.
type CommandHandlers struct {
	FailsCommandHandler        *command.FailsCommandHandler
	SubscriptionCommandHandler *command.SubscriptionCommandHandler
	UserCommandHandler         *command.UserCommandHandler
	GameCommandHandler         *command.GameCommandHandler
}

// QueryHandlers holds all query handlers.
type QueryHandlers struct {
	FailsQueryHandler        *query.FailQueryHandler
	UserQueryHandler         *query.UserQueryHandler
	SubscriptionQueryHandler *query.SubscriptionQueryHandler
	GameQueryHandler         *query.GameQueryHandler
}

// Repositories holds all repositories.
type Repositories struct {
	FailRepository         *repository.FailRepository
	UserRepository         *repository.UserRepository
	SubscriptionRepository *repository.SubscriptionRepository
	GameRepository         *repository.GameRepository
}

// AppInit initializes the application with the provided database connection.
func AppInit(dataBaseConnection *sql.DB) *App {
	databaseImpl := &repository.DatabaseImpl{DB: dataBaseConnection}
	services := createServices()
	repositories := createRepositories(databaseImpl)
	commandHandlers := createCommandHandlers(repositories, services)
	queryHandlers := createQueryHandlers(repositories)
	jobs := createJobs(queryHandlers, commandHandlers)
	endpoints := createEndpoints(commandHandlers, queryHandlers, services)

	return &App{
		DatabaseImpl:    databaseImpl,
		Repositories:    repositories,
		CommandHandlers: commandHandlers,
		QueryHandlers:   queryHandlers,
		Jobs:            jobs,
		Endpoints:       endpoints,
		Services:        services,
	}
}

func createRepositories(databaseImpl *repository.DatabaseImpl) *Repositories {
	return &Repositories{
		SubscriptionRepository: &repository.SubscriptionRepository{DatabaseImpl: databaseImpl},
		UserRepository:         &repository.UserRepository{DatabaseImpl: databaseImpl},
		FailRepository:         &repository.FailRepository{DatabaseImpl: databaseImpl},
		GameRepository:         &repository.GameRepository{DatabaseImpl: databaseImpl},
	}
}

func createCommandHandlers(repositories *Repositories, services *Services) *CommandHandlers {
	failsHandler := &command.FailsCommandHandler{FailRepository: repositories.FailRepository}
	gameHandler := &command.GameCommandHandler{GameRepository: repositories.GameRepository}
	subscriptionCommandHandler := &command.SubscriptionCommandHandler{
		SteamService:           services.SteamService,
		SubscriptionRepository: repositories.SubscriptionRepository,
		FailsCommandHandler:    failsHandler,
		GameCommandHandler:     gameHandler,
	}

	return &CommandHandlers{
		FailsCommandHandler:        failsHandler,
		SubscriptionCommandHandler: subscriptionCommandHandler,
		UserCommandHandler:         &command.UserCommandHandler{UserRepository: repositories.UserRepository},
		GameCommandHandler:         &command.GameCommandHandler{GameRepository: repositories.GameRepository},
	}
}

func createQueryHandlers(repositories *Repositories) *QueryHandlers {
	return &QueryHandlers{
		FailsQueryHandler:        &query.FailQueryHandler{FailRepository: repositories.FailRepository},
		UserQueryHandler:         &query.UserQueryHandler{UserRepository: repositories.UserRepository},
		SubscriptionQueryHandler: &query.SubscriptionQueryHandler{SubscriptionRepository: repositories.SubscriptionRepository},
		GameQueryHandler:         &query.GameQueryHandler{GameRepository: repositories.GameRepository},
	}
}

func createJobs(queryHandlers *QueryHandlers, commandHandlers *CommandHandlers) *cron.Jobs {

	return &cron.Jobs{
		FailCron: &cron.FailCron{
			SubscriptionCommandHandler: commandHandlers.SubscriptionCommandHandler,
			SubscriptionQueryHandler:   queryHandlers.SubscriptionQueryHandler,
		},
	}
}

func createServices() *Services {
	return &Services{
		SteamService: &service.SteamService{},
		AuthService:  &service.AuthService{},
	}
}

func createEndpoints(commandHandlers *CommandHandlers, queryHandlers *QueryHandlers, services *Services) *endpoint.Endpoints {
	router := mux.NewRouter()

	return &endpoint.Endpoints{
		Router: router,
		UserEndpoint: &endpoint.UserEndpoint{
			Router:                     router,
			SteamService:               services.SteamService,
			UserCommandHandler:         commandHandlers.UserCommandHandler,
			UserQueryHandler:           queryHandlers.UserQueryHandler,
			SubscriptionCommandHandler: commandHandlers.SubscriptionCommandHandler,
			AuthService:                services.AuthService,
		},
		FailEndpoint: &endpoint.FailEndpoint{
			Router:           router,
			FailQueryHandler: queryHandlers.FailsQueryHandler,
		},
		SubscriptionEndpoint: &endpoint.SubscriptionEndpoint{
			Router:                     router,
			SubscriptionCommandHandler: commandHandlers.SubscriptionCommandHandler,
			SubscriptionQueryHandler:   queryHandlers.SubscriptionQueryHandler,
			AuthService:                services.AuthService,
		},
		GamesEndpoint: &endpoint.GamesEndpoint{
			Router:           router,
			GameQueryHandler: queryHandlers.GameQueryHandler,
		},
	}
}
