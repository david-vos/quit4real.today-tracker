package endpoint

import (
	"github.com/gorilla/mux"
	"net/http"
	"quit4real.today/logger"
)

type Endpoints struct {
	// This router pointer is global, it shouldn't be. I have yet to think of a clean way to instantiate every
	// endpoint on the same router while still keeping all the Endpoints separated from each other :thinking:
	// I think mux has a router.sub-route but this will have to work for now...
	Router               *mux.Router
	UserEndpoint         *UserEndpoint
	FailEndpoint         *FailEndpoint
	SubscriptionEndpoint *SubscriptionEndpoint
}

func Init(endpoints *Endpoints) {
	// Start the HTTP server in a separate goroutine
	endpoints.UserEndpoint.User()
	endpoints.FailEndpoint.Fail()
	endpoints.SubscriptionEndpoint.Subscription()

	go func() {
		logger.Info("Starting the server. If you do not receive and error it has started successfully")
		if err := http.ListenAndServe(":8080", endpoints.Router); err != nil {
			logger.Fail("Failed to start server: " + err.Error())
		}
	}()
}
