package client

import (
	"github.com/thelotter-enterprise/usergo/shared"
)

type ServiceMiddleware func(UserService) UserService

type UserService interface {
	GetUserByID(id int) shared.HTTPResponse
}
