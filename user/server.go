package user

import (
	"context"
	"net/http"
	"os"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//NewServer ...
func NewServer(ctx context.Context, endpoints Endpoints) http.Handler {
	router := mux.NewRouter()

	getUserByIDHandler := httptransport.NewServer(
		endpoints.GetUserByID,
		decodeGetUserByIDRequest,
		encodeReponseToJSON,
	)

	router.Methods("GET").Path("/user/{id}").Handler(getUserByIDHandler)

	return handlers.LoggingHandler(os.Stdout, router)
}
