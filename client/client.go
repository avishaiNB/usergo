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
func (client *ServiceClient) GetUserByID(ctx context.Context, id int) (om.ByIDResponse, error) {
	ep, _ := NewUserByIDEndpoint(id)
	res, _ := ep(ctx, id)
	response := res.(om.ByIDResponse)
	return response, nil
}
