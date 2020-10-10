package core

import (
	"context"
	"net/http"
	"time"

	httpkit "github.com/go-kit/kit/transport/http"
)

// CorrelationIDHeaderKey ..
type CorrelationIDHeaderKey string

// TimeoutHeaderKey ..
type TimeoutHeaderKey string

const (
	// MaxTimeout is 15 seconds
	MaxTimeout time.Duration = time.Second * 15

	// CorrelationIDKey ...
	CorrelationIDKey CorrelationIDHeaderKey = "correlation_id"

	// TimeoutKey ..
	TimeoutKey TimeoutHeaderKey = "timeout"
)

// Ctx ...
type Ctx struct {
	Context context.Context
}

// NewCtx will create a new context
func NewCtx(ctx context.Context) Ctx {
	return Ctx{
		Context: ctx,
	}
}

// CalcTimeout will return the timeout (deadline) for waiting an external response to come back
// TODO: now I return max, need to change it
func (ctx Ctx) CalcTimeout() time.Duration {
	return MaxTimeout
}

//func AddToCtx

// ReadCtx ...
func ReadCtx(ctx context.Context, r *http.Request) context.Context {
	correlationid := r.Header.Get(string(CorrelationIDKey))
	timeout := r.Header.Get(string(TimeoutKey))

	ctx = context.WithValue(ctx, CorrelationIDKey, correlationid)
	ctx = context.WithValue(ctx, TimeoutKey, timeout)

	return ctx
}

// ReadCtxBefore ...
func ReadCtxBefore() httpkit.ServerOption {
	return httpkit.ServerBefore(ReadCtx)
}
