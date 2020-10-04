package client

import (
	"context"

	"github.com/thelotter-enterprise/usergo/shared"
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
func (client *ServiceClient) GetUserByID(ctx context.Context, id int) (shared.ByIDResponse, error) {
	ep, _ := NewUserByIDEndpoint(id)
	res, _ := ep(ctx, id)
	response := res.(shared.ByIDResponse)
	return response, nil
}
