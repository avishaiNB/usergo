package http

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/go-kit/kit/endpoint"
	kitlogger "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	tlectx "github.com/thelotter-enterprise/usergo/core/context"
	tlehttp "github.com/thelotter-enterprise/usergo/core/transports/http"
	"github.com/thelotter-enterprise/usergo/core/utils"
	"github.com/thelotter-enterprise/usergo/shared"
	"github.com/thelotter-enterprise/usergo/svc"
	"github.com/thelotter-enterprise/usergo/svc/transport"
)

// NewService ..
func NewService(svcEndpoints transport.Endpoints, options []kithttp.ServerOption, logger kitlogger.Logger) http.Handler {
	// set-up router and initialize http endpoints
	var (
		router        = mux.NewRouter()
		errorLogger   = kithttp.ServerErrorLogger(logger)
		errorEncoder  = kithttp.ServerErrorEncoder(encodeErrorResponse)
		contextReader = tlectx.ReadBefore()
	)

	options = append(options, errorLogger, errorEncoder, contextReader)

	getUserByIDHandler := kithttp.NewServer(
		svcEndpoints.UserByIDEndpoint,
		decodeUserByIDRequest,
		encodeUserByIDReponse,
		options...)

	router.Methods("GET").Path(shared.UserByIDServerRoute).Handler(getUserByIDHandler)

	// logger.Info(ctx, fmt.Sprintf("adding route http://%s/%s", server.Address, endpoint.Path))

	return handlers.LoggingHandler(os.Stdout, router)

	//logger.Info(ctx, fmt.Sprintf("http server started and listen on %s", server.Address))
	//http.ListenAndServe(server.Address, server.Handler)
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

func encodeUserByIDReponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func decodeUserByIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req interface{}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	return req, err
}

func encodeErrorResponse(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	default:
		return http.StatusInternalServerError
	}
}
