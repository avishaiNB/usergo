package client

import (
	"os"

	"github.com/go-kit/kit/log"
	"github.com/thelotter-enterprise/usergo/shared"
)

// UserServiceClientMiddleware used to chain behaviors on the UserService using middleware pattern
type UserServiceClientMiddleware func(UserServiceClient) UserServiceClient

// UserServiceClient defines all the APIs available for the service
type UserServiceClient interface {
	// Gets the user by an ID
	GetUserByID(id int) shared.HTTPResponse

	// Gets the user by email
	GetUserByEmail(email string) shared.HTTPResponse
}

// ServiceClient is a facade for all APIs exposed by the service
type ServiceClient struct {
	Logger      log.Logger
	ServiceName string
	SD          ServiceDiscovery
}

// NewServiceClient will create a new instance of ServiceClient
func NewServiceClient() ServiceClient {
	consulAddress := "localhost:8080"

	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "listen", ":8080", "caller", log.DefaultCaller)

	sd, _ := NewServiceDiscovery(logger, consulAddress)

	client := ServiceClient{
		Logger:      logger,
		ServiceName: "user",
		SD:          sd,
	}
	return client
}
