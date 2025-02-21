package src

import (
	"database/sql"
	"github.com/gorilla/mux"
	"quit4real.today/src/cronJobs"
	"quit4real.today/src/endpoint"
	"quit4real.today/src/handler/command"
	"quit4real.today/src/handler/query"
	"quit4real.today/src/handler/service"
	"quit4real.today/src/handler/service/impl"
	repoImpl "quit4real.today/src/repository/impl"
)

// App represents the main application structure.
type App struct {
	DatabaseImpl    *repoImpl.DatabaseImpl
	Endpoints       *endpoint.Endpoints
	CommandHandlers *CommandHandlers
	QueryHandlers   *QueryHandlers
	Repositories    *Repositories
	Jobs            *cronJobs.Jobs
	Services        *Services
}

// Services holds service dependencies as interfaces.
type Services struct {
	SteamService service.SteamService
	AuthService  service.AuthService
}

// CommandHandlers holds all command handlers.
type CommandHandlers struct {
	FailsCommandHandler        *command.FailsCommandHandlerImpl
	SubscriptionCommandHandler *command.SubscriptionCommandHandlerImpl
	UserCommandHandler         *command.UserCommandHandlerImpl
	GameCommandHandler         *command.GameCommandHandlerImpl
}

// QueryHandlers holds all query handlers.
type QueryHandlers struct {
	FailsQueryHandler        *query.FailsQueryHandlerImpl
	UserQueryHandler         *query.UserQueryHandlerImpl
	SubscriptionQueryHandler *query.SubscriptionQueryHandlerImpl
	GameQueryHandler         *query.GameQueryHandlerImpl
}

// Repositories holds all repositories.
type Repositories struct {
	FailRepository         *repoImpl.FailRepositoryImpl
	UserRepository         *repoImpl.UserRepositoryImpl
	SubscriptionRepository *repoImpl.SubscriptionRepositoryImpl
	GameRepository         *repoImpl.GameRepositoryImpl
}

// AppInit initializes the application with the provided database connection and services.
func AppInit(dataBaseConnection *sql.DB, steamService service.SteamService, authService service.AuthService) *App {
	databaseImpl := &repoImpl.DatabaseImpl{DB: dataBaseConnection}
	services := createServices(steamService, authService)
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

func createServices(steamService service.SteamService, authService service.AuthService) *Services {
	return &Services{
		SteamService: steamService,
		AuthService:  authService,
	}
}

func createRepositories(databaseImpl *impl.DatabaseImpl) *Repositories {
	return &Repositories{
		SubscriptionRepository: &repoImpl.SubscriptionRepositoryImpl{DatabaseImpl: databaseImpl},
		UserRepository:         &repoImpl.UserRepositoryImpl{DatabaseImpl: databaseImpl},
		FailRepository:         &repoImpl.FailRepositoryImpl{DatabaseImpl: databaseImpl},
		GameRepository:         &repoImpl.GameRepositoryImpl{DatabaseImpl: databaseImpl},
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
		UserCommandHandler:         &command.UserCommandHandlerImpl{UserRepository: repositories.UserRepository},
		GameCommandHandler:         gameHandler,
	}
}

func createQueryHandlers(repositories *Repositories) *QueryHandlers {
	return &QueryHandlers{
		FailsQueryHandler:        &query.FailsQueryHandlerImpl{FailRepository: repositories.FailRepository},
		UserQueryHandler:         &query.UserQueryHandlerImpl{UserRepository: repositories.UserRepository},
		SubscriptionQueryHandler: &query.SubscriptionQueryHandler{SubscriptionRepository: repositories.SubscriptionRepository},
		GameQueryHandler:         &query.GameQueryHandlerImpl{GameRepository: repositories.GameRepository},
	}
}

func createJobs(queryHandlers *QueryHandlers, commandHandlers *CommandHandlers) *cronJobs.Jobs {
	// Assuming cronJobs.Jobs depends on command and query handlers
	return &cronJobs.Jobs{}
}

func createEndpoints(commandHandlers *CommandHandlers, queryHandlers *QueryHandlers, services *Services) *endpoint.Endpoints {
	router := mux.NewRouter()
	return &endpoint.Endpoints{
		Router: router,
		SubscriptionEndpoint: &endpoint.SubscriptionEndpoint{
			Router:                     router,
			SubscriptionCommandHandler: commandHandlers.SubscriptionCommandHandler,
			SubscriptionQueryHandler:   queryHandlers.SubscriptionQueryHandler,
			AuthService:                services.AuthService,
		},
	}
}
