package userclient

import (
	om "github.com/thelotter-enterprise/usergo/usershared"
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
func (client *ServiceClient) GetUserByID(id int) (om.User, error) {
	ep, _ := NewUserByIDEndpoint(id)
	println(ep)
	var user om.User
	return user, nil
}
