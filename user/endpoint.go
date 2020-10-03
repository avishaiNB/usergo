package user

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

// Endpoints ...
type Endpoints struct {
	GetUserByID endpoint.Endpoint
}

// MakeEndpoints ...
func MakeEndpoints(s Service) Endpoints {
	return Endpoints{
		GetUserByID: makeGetUserByIDEndpoint(s),
	}
}

func makeGetUserByIDEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ByIDRequest)
		user, err := s.GetUserByID(ctx, req.ID)
		return NewGetUserResponse(user), err
	}
}
