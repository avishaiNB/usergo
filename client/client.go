package client

import (
	"context"

	om "github.com/thelotter-enterprise/usergo/shared"
)

// ServiceClient ...
type ServiceClient struct {
}

// NewServiceClient ...
func NewServiceClient() (ServiceClient, error) {
	client := ServiceClient{}
	return client, nil
}

// GetUserByID ..
// TODO: we should not return the User we need to have some kind of a wrapper
func (client *ServiceClient) GetUserByID(ctx context.Context, id int) (om.User, error) {
	ep, _ := NewUserByIDEndpoint(id)
	println(ep)
	var user om.User
	return user, nil
}
