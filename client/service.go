package client

import (
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
}

// NewServiceClient will create a new instance of ServiceClient
func NewServiceClient() ServiceClient {
	client := ServiceClient{}
	return client
}
