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
func (client *ServiceClient) GetUserByID(ctx context.Context, id int) shared.HTTPResponse {
	serviceClient := serviceendpoints.NewUserByIDServiceClient(client.Router)
	serviceClient.WithContext(ctx)
	serviceClient.WithParams(map[string]interface{}{"ID": id})
	serviceClient.WithCircuitBreaker("get-user-by-id", shared.NewHystrixCommandConfig())
	serviceClient.BuildEndpoints()
	serviceClient.Exec()
	return serviceClient.GetResult()
}

// GetUserByID2 ..
func (client *ServiceClient) GetUserByID2(ctx context.Context, id int) shared.HTTPResponse {
	var svc UserService
	var endpoints []ProxyEndpoint
	endpoints = append(endpoints, makeProxyEndpoint(id))
	svc = proxyingMiddleware(context.Background(), endpoints)(svc)
	return svc.GetUserByID(id)
}
