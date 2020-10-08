package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/thelotter-enterprise/usergo/core"
	"github.com/thelotter-enterprise/usergo/shared"
)

type proxymw struct {
	breakermw  endpoint.Middleware
	limmitermw endpoint.Middleware
	router     *mux.Router
}

func newProxyMiddleware(breakermw endpoint.Middleware, limittermw endpoint.Middleware, router *mux.Router) proxymw {
	return proxymw{
		breakermw:  breakermw,
		limmitermw: limittermw,
		router:     router,
	}
}

func (mw proxymw) userByIDMiddleware(ctx context.Context, id int) UserServiceMiddleware {
	var endpointer sd.FixedEndpointer

	tgt, _ := mw.router.Schemes("http").Host("localhost:8080").Path(shared.UserByIDRoute).URL("id", strconv.Itoa(id))
	e := httptransport.NewClient("GET", tgt, core.EncodeRequestToJSON, decodeGetUserByIDResponse).Endpoint()
	e = mw.breakermw(e)
	e = mw.limmitermw(e)
	endpointer = append(endpointer, e)

	lb := core.NewLoadBalancer(endpointer)
	retry := lb.DefaultRoundRobinWithRetryEndpoint(ctx)

	return func(next UserService) UserService {
		out := core.ProxyMiddlewareData{Context: ctx, Next: next, This: retry}

		return userByIDProxyMiddleware{
			mw: out,
		}
	}
}

func decodeGetUserByIDResponse(_ context.Context, r *http.Response) (interface{}, error) {
	if r.StatusCode != http.StatusOK {
		return nil, errors.New(r.Status)
	}
	var resp shared.ByIDResponse
	err := json.NewDecoder(r.Body).Decode(&resp)
	return resp, err
}

type userByIDProxyMiddleware struct {
	mw core.ProxyMiddlewareData
}

// GetUserByID will execute the endpoint using the middleware and will constract an shared.HTTPResponse
func (proxymw userByIDProxyMiddleware) GetUserByID(id int) core.HTTPResponse {
	var res interface{}
	var err error
	circuitOpen := false
	statusCode := 200

	if res, err = proxymw.mw.This(proxymw.mw.Context, id); err != nil {
		// TODO: need a refactor to analyze the response
		circuitOpen = true
		statusCode = 500
	}

	return core.HTTPResponse{
		Result:        res,
		Error:         err,
		CircuitOpened: circuitOpen,
		StatusCode:    statusCode,
	}
}

// GetUserByEmail will proxy the implementation to the responsible middleware
// We do this to satisfy the service interface
func (proxymw userByIDProxyMiddleware) GetUserByEmail(email string) core.HTTPResponse {
	svc := proxymw.mw.Next.(UserService)
	return svc.GetUserByEmail(email)
}
