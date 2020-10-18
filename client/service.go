package client

import tlehttp "github.com/thelotter-enterprise/usergo/core/transports/http"

// ServiceMiddleware used to chain behaviors on the UserService using middleware pattern
type ServiceMiddleware func(Service) Service

// Service defines all the APIs available for the service
type Service interface {
	// Gets the user by an ID
	GetUserByID(id int) tlehttp.Response

	// Gets the user by email
	GetUserByEmail(email string) tlehttp.Response
}
