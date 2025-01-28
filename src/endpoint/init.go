package endpoint

import (
	"fmt"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"quit4real.today/logger"
)

type Endpoints struct {
	Router               *mux.Router
	UserEndpoint         *UserEndpoint
	FailEndpoint         *FailEndpoint
	SubscriptionEndpoint *SubscriptionEndpoint
	GamesEndpoint        *GamesEndpoint
}

// getLocalIPs returns a list of local IPv4 addresses
func getLocalIPs() []string {
	var ips []string
	interfaces, err := net.Interfaces()
	if err != nil {
		logger.Fail("Error fetching network interfaces: " + err.Error())
		return ips
	}

	for _, iface := range interfaces {
		// Skip interfaces that are down or loopback interfaces
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			logger.Fail("Error fetching addresses for interface " + iface.Name + ": " + err.Error())
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// Only include IPv4 addresses
			if ip != nil && ip.To4() != nil {
				ips = append(ips, ip.String())
			}
		}
	}

	return ips
}

func Init(endpoints *Endpoints) {
	// Start the HTTP server in a separate goroutine
	endpoints.UserEndpoint.User()
	endpoints.FailEndpoint.Fail()
	endpoints.SubscriptionEndpoint.Subscription()
	endpoints.GamesEndpoint.Games()

	go func() {
		listener, err := net.Listen("tcp", "0.0.0.0:8080")
		if err != nil {
			logger.Fail("Failed to create listener: " + err.Error())
			return
		}

		// Log the available IPs with the port
		ips := getLocalIPs()
		if len(ips) > 0 {
			for _, ip := range ips {
				logger.Info(fmt.Sprintf("Server is accessible at http://%s:8080", ip))
			}
		} else {
			logger.Debug("No active network interfaces found. Server might not be reachable.")
		}

		// Start the server using the listener
		if err := http.Serve(listener, endpoints.Router); err != nil {
			logger.Fail("Failed to start server: " + err.Error())
		}
	}()
}
