package client

import (
	"context"

	"github.com/gorilla/mux"
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
	ep := newUserByIDEndpoint(ctx, id, client.Router)
	ep.build()
	ep.exec()
	return ep.result()
}
