package client

import (
	"context"

	"github.com/gorilla/mux"
	"github.com/thelotter-enterprise/usergo/client/serviceendpoints"
	"github.com/thelotter-enterprise/usergo/shared"
)

// ServiceClient ...
type ServiceClient struct {
	Router *mux.Router
}

// NewServiceClient ...
func NewServiceClient() (ServiceClient, error) {
	client := ServiceClient{
		Router: mux.NewRouter(),
	}
	return client, nil
}

// GetUserByID ..
func (client *ServiceClient) GetUserByID(ctx context.Context, id int) (shared.ByIDResponse, error) {
	serviceClient := serviceendpoints.NewUserByIDServiceClient(ctx, id, client.Router)
	serviceClient.BuildEndpoints()
	serviceClient.Exec()
	return serviceClient.GetResult()
}
