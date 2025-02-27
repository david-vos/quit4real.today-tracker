package src

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/yohcop/openid-go"
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
	SteamService        service.SteamService
	AuthService         service.AuthService
	TrackerService      service.TrackerService
	SubscriptionService service.SubscriptionService
	UserService         service.UserService
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
	FailQueryHandler         *query.FailQueryHandlerImpl
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
	TrackerRepository      *repoImpl.TrackerRepositoryImpl
}

// AppInit initializes the application with all required services and components.
func AppInit(databaseConnection *sql.DB) *App {
	// Initialize the database implementation
	databaseImpl := &repoImpl.DatabaseImpl{DB: databaseConnection}

	// Initialize repositories
	repositories := createRepositories(databaseImpl)

	// Initialize services
	services := createServices(repositories)

	// Initialize command and query handlers
	commandHandlers := createCommandHandlers(repositories, services)
	queryHandlers := createQueryHandlers(repositories)

	// Fill the subscriptionService
	services.SubscriptionService = impl.NewSubscriptionServiceImpl(*queryHandlers.SubscriptionQueryHandler, services.SteamService)

	// Initialize cron jobs
	jobs := createJobs(queryHandlers, commandHandlers, services)

	// Initialize endpoints
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

// createServices initializes and returns all application services.
func createServices(repositories *Repositories) *Services {
	// THIS repositories IS REALLY BAD BUT I NEED SLEEP!
	steamService := impl.NewSteamServiceImpl()
	authService := impl.NewAuthServiceImpl(*openid.NewOpenID(nil)) // Use default HTTP client for OpenID
	trackerService := impl.NewTrackerServiceImpl(repositories.TrackerRepository)
	userService := impl.NewUserServiceImpl(trackerService, steamService)

	return &Services{
		SteamService:   steamService,
		AuthService:    authService,
		TrackerService: trackerService,
		UserService:    userService,
	}
}

// createRepositories initializes and returns all repositories.
func createRepositories(databaseImpl *repoImpl.DatabaseImpl) *Repositories {
	return &Repositories{
		FailRepository:         &repoImpl.FailRepositoryImpl{DatabaseImpl: databaseImpl},
		UserRepository:         &repoImpl.UserRepositoryImpl{DatabaseImpl: databaseImpl},
		SubscriptionRepository: &repoImpl.SubscriptionRepositoryImpl{DatabaseImpl: databaseImpl},
		GameRepository:         &repoImpl.GameRepositoryImpl{DatabaseImpl: databaseImpl},
		TrackerRepository:      &repoImpl.TrackerRepositoryImpl{DatabaseImpl: databaseImpl},
	}
}

// createCommandHandlers initializes and returns all command handlers.
func createCommandHandlers(repositories *Repositories, services *Services) *CommandHandlers {
	return &CommandHandlers{
		FailsCommandHandler: &command.FailsCommandHandlerImpl{
			FailRepository: repositories.FailRepository,
		},
		SubscriptionCommandHandler: &command.SubscriptionCommandHandlerImpl{
			SteamService:           services.SteamService,
			SubscriptionRepository: repositories.SubscriptionRepository,
			FailsCommandHandler:    &command.FailsCommandHandlerImpl{FailRepository: repositories.FailRepository},
			GameCommandHandler:     &command.GameCommandHandlerImpl{GameRepository: repositories.GameRepository},
		},
		UserCommandHandler: &command.UserCommandHandlerImpl{
			UserRepository: repositories.UserRepository,
		},
		GameCommandHandler: &command.GameCommandHandlerImpl{
			GameRepository: repositories.GameRepository,
		},
	}
}

// createQueryHandlers initializes and returns all query handlers.
func createQueryHandlers(repositories *Repositories) *QueryHandlers {
	return &QueryHandlers{
		FailQueryHandler:         &query.FailQueryHandlerImpl{FailRepository: repositories.FailRepository},
		UserQueryHandler:         &query.UserQueryHandlerImpl{UserRepository: repositories.UserRepository},
		SubscriptionQueryHandler: &query.SubscriptionQueryHandlerImpl{SubscriptionRepository: repositories.SubscriptionRepository},
		GameQueryHandler:         &query.GameQueryHandlerImpl{GameRepository: repositories.GameRepository},
	}
}

// createJobs initializes and returns cron jobs.
func createJobs(queryHandlers *QueryHandlers, commandHandlers *CommandHandlers, services *Services) *cronJobs.Jobs {
	return &cronJobs.Jobs{
		FailCron: &cronJobs.FailCronImpl{
			SubscriptionQueryHandler:   queryHandlers.SubscriptionQueryHandler,
			SubscriptionCommandHandler: commandHandlers.SubscriptionCommandHandler,
			UserQueryHandler:           queryHandlers.UserQueryHandler,
			SteamService:               services.SteamService,
			SubscriptionService:        services.SubscriptionService,
			TrackerService:             services.TrackerService,
		},
	}
}

// createEndpoints initializes and returns all endpoints.
func createEndpoints(commandHandlers *CommandHandlers, queryHandlers *QueryHandlers, services *Services) *endpoint.Endpoints {
	router := mux.NewRouter()

	return &endpoint.Endpoints{
		Router: router,
		UserEndpoint: &endpoint.UserEndpoint{
			Router:                     router,
			UserQueryHandler:           queryHandlers.UserQueryHandler,
			UserCommandHandler:         commandHandlers.UserCommandHandler,
			SubscriptionCommandHandler: commandHandlers.SubscriptionCommandHandler,
			SteamService:               services.SteamService,
			AuthService:                services.AuthService,
			UserService:                services.UserService,
		},
		FailEndpoint: &endpoint.FailEndpoint{
			Router:           router,
			FailQueryHandler: queryHandlers.FailQueryHandler,
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
