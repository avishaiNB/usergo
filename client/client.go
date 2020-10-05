package client

import (
	"context"

	"github.com/thelotter-enterprise/usergo/shared"
)

// GetUserByID , if found will return shared.HTTPResponse containing the user requested information
// If an error occurs it will hold error information that cab be used to decide how to proceed
func (client *ServiceClient) GetUserByID(ctx context.Context, id int) shared.HTTPResponse {
	var svc UserServiceClient
	commandName := "get_user_by_id"

	endpoints := makeEndpoints(id)
	input := shared.MakeProxyMiddlewareData(ctx, commandName, endpoints)

	svc = makeProxyMiddleware(input)(svc)
	svc = makeLoggingMiddleware(client.Logger)(svc)
	svc = makeInstrumentingMiddleware(client.Name, commandName)(svc)
	return svc.GetUserByID(id)
}
