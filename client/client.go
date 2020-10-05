package client

import (
	"context"

	"github.com/thelotter-enterprise/usergo/shared"
)

// GetUserByID , if found will return shared.HTTPResponse containing the user requested information
// If an error occurs it will hold error information that cab be used to decide how to proceed
func (client *ServiceClient) GetUserByID(ctx context.Context, id int) shared.HTTPResponse {
	var svc UserServiceClient
	commandName := "GetUserByID"

	endpoints := makeUserByIDEndpoints(id)
	input := shared.MakeDefaultMiddlewareInput(ctx, commandName, endpoints)
	proxyMiddleware := makeUserByIDMiddleware(input)
	svc = proxyMiddleware(svc)
	svc = makeLoggingMiddleware(client.Logger)(svc)
	return svc.GetUserByID(id)
}
