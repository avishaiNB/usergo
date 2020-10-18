package svc

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/thelotter-enterprise/usergo/core"
	tlehttp "github.com/thelotter-enterprise/usergo/core/http"
	tletracer "github.com/thelotter-enterprise/usergo/core/tracer"
	"github.com/thelotter-enterprise/usergo/shared"
)

// UserHTTPEndpoints ...
type UserHTTPEndpoints struct {
	HTTPEndpoints *tlehttp.HTTPEndpoints
	Service       Service
	Log           core.Log
	Tracer        tletracer.Tracer
}

// NewUserHTTPEndpoints ...
func NewUserHTTPEndpoints(log core.Log, tracer tletracer.Tracer, service Service) *UserHTTPEndpoints {
	userEndpoints := UserHTTPEndpoints{
		Log:           log,
		Tracer:        tracer,
		Service:       service,
		HTTPEndpoints: &tlehttp.HTTPEndpoints{},
	}

	userEndpoints.HTTPEndpoints = userEndpoints.makeEndpoints()

	return &userEndpoints
}

func (ue UserHTTPEndpoints) makeEndpoints() *tlehttp.HTTPEndpoints {
	var endpoints tlehttp.HTTPEndpoints
	var serverEndpoints []tlehttp.HTTPEndpoint

	userbyid := tlehttp.HTTPEndpoint{
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

		decoder := core.NewDecoder()

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
		ue.Log.Logger.Log("method", "EncodeReponseToJSONFunc", "error", err)
	}
	return err
}

func (ue UserHTTPEndpoints) decodeUserByIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		ue.Log.Logger.Log(
			"level", "error",
			"method", "DecodeRequestFromJSONFunc",
			"error", err,
		)
	}

	return req, err
}
