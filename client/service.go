package client

import tlehttp "github.com/thelotter-enterprise/usergo/core/http"

// UserServiceMiddleware used to chain behaviors on the UserService using middleware pattern
type UserServiceMiddleware func(UserService) UserService

// UserService defines all the APIs available for the service
type UserService interface {
	// Gets the user by an ID
	GetUserByID(id int) tlehttp.Response

	// Gets the user by email
	GetUserByEmail(email string) tlehttp.Response
}
