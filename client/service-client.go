package client

import (
	"context"

	"github.com/thelotter-enterprise/usergo/shared"
)

// ServiceClient ...
type ServiceClient struct {
}

// NewServiceClient ...
func NewServiceClient() ServiceClient {
	client := ServiceClient{}
	return client
}

// GetUserByID ..
func (client *ServiceClient) GetUserByID(ctx context.Context, id int) shared.HTTPResponse {
	var svc UserService
	var endpoints []ProxyEndpoint
	endpoints = append(endpoints, makeProxyEndpoint(id))
	input := NewMiddlewareInput(context.Background(), "get-user-by-id", endpoints)
	svc = proxyingMiddleware(input)(svc)
	return svc.GetUserByID(id)
}
