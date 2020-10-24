package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	tlehttp "github.com/thelotter-enterprise/usergo/core/transports/http"
	"github.com/thelotter-enterprise/usergo/core/utils"
	"github.com/thelotter-enterprise/usergo/shared"
	"github.com/thelotter-enterprise/usergo/svc"
)

// Endpoints holds all Go kit endpoints for the Order service.
type Endpoints struct {
	UserByIDEndpoint endpoint.Endpoint
}

// MakeEndpoints initializes all Go kit endpoints for the Order service.
func MakeEndpoints(s svc.Service) Endpoints {
	return Endpoints{
		UserByIDEndpoint: makeUserByIDEndpoint(s),
	}
}

func makeUserByIDEndpoint(service svc.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var err error
		var req tlehttp.Request
		var data shared.ByIDRequestData

		decoder := utils.NewDecoder()

		err = decoder.MapDecode(request, &req)
		err = decoder.MapDecode(req.Data, &data)
		req.Data = data

		user, err := service.GetUserByID(ctx, data.ID)
		return shared.NewUserResponse(user), err
	}
}
