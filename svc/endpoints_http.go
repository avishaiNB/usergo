package svc

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	tletracer "github.com/thelotter-enterprise/usergo/core/tracer"
	tlehttp "github.com/thelotter-enterprise/usergo/core/transports/http"
	"github.com/thelotter-enterprise/usergo/core/utils"
	"github.com/thelotter-enterprise/usergo/shared"

	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
)

// UserHTTPEndpoints ...
type UserHTTPEndpoints struct {
	HTTPEndpoints *tlehttp.Endpoints
	Service       Service
	Logger        *tlelogger.Manager
	Tracer        tletracer.Tracer
}

// NewUserHTTPEndpoints ...
func NewUserHTTPEndpoints(logger *tlelogger.Manager, tracer tletracer.Tracer, service Service) *UserHTTPEndpoints {
	userEndpoints := UserHTTPEndpoints{
		Logger:        logger,
		Tracer:        tracer,
		Service:       service,
		HTTPEndpoints: &tlehttp.Endpoints{},
	}

	userEndpoints.HTTPEndpoints = userEndpoints.makeEndpoints()

	return &userEndpoints
}

func (ue UserHTTPEndpoints) makeEndpoints() *tlehttp.Endpoints {
	var endpoints tlehttp.Endpoints
	var serverEndpoints []tlehttp.Endpoint

	userbyid := tlehttp.Endpoint{
		Endpoint: makeUserByIDEndpoint(ue.Service),
		Enc:      ue.encodeUserByIDReponse,
		Dec:      ue.decodeUserByIDRequest,
		Method:   "GET",
		Path:     shared.UserByIDServerRoute,
	}

	serverEndpoints = append(serverEndpoints, userbyid)
	endpoints.ServerEndpoints = serverEndpoints
	return &endpoints
}

func makeUserByIDEndpoint(service Service) endpoint.Endpoint {
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

func (ue UserHTTPEndpoints) encodeUserByIDReponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		logger := *ue.Logger
		logger.Error(
			ctx,
			"encodeUserByIDReponse",
			"method", "EncodeReponseToJSONFunc", "error", err)
	}
	return err
}

func (ue UserHTTPEndpoints) decodeUserByIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		logger := *ue.Logger
		logger.Error(
			ctx,
			"decodeUserByIDRequest",
			"level", "error",
			"method", "DecodeRequestFromJSONFunc",
			"error", err,
		)
	}

	return req, err
}
