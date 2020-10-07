package client

import "github.com/thelotter-enterprise/usergo/core"

// UserServiceMiddleware used to chain behaviors on the UserService using middleware pattern
type UserServiceMiddleware func(UserService) UserService

// UserService defines all the APIs available for the service
type UserService interface {
	// Gets the user by an ID
	GetUserByID(id int) core.HTTPResponse

	// Gets the user by email
	GetUserByEmail(email string) core.HTTPResponse
}
