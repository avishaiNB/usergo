package client

import (
	"context"

	"github.com/go-kit/kit/log"
	"github.com/thelotter-enterprise/usergo/core"
)

// ServiceClient is a facade for all APIs exposed by the service
type ServiceClient struct {
	Logger      log.Logger
	SD          *core.ServiceDiscoverator
	ServiceName string
}

// NewServiceClient will create a new instance of ServiceClient
func NewServiceClient(logger log.Logger, sd *core.ServiceDiscoverator, serviceName string) ServiceClient {
	client := ServiceClient{
		Logger:      logger,
		SD:          sd,
		ServiceName: serviceName,
	}
	return client
}

// GetUserByID , if found will return shared.HTTPResponse containing the user requested information
// If an error occurs it will hold error information that cab be used to decide how to proceed
func (client *ServiceClient) GetUserByID(ctx context.Context, id int) core.HTTPResponse {
	var svc UserService
	commandName := "get_user_by_id"

	endpoints := makeEndpoints(id)
	input := core.MakeProxyMiddlewareData(ctx, commandName, endpoints)

	svc = makeProxyMiddleware(input)(svc)
	svc = makeLoggingMiddleware(client.Logger)(svc)
	svc = makeInstrumentingMiddleware(client.ServiceName, commandName)(svc)
	return svc.GetUserByID(id)
}
