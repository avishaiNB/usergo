package core

import (
	"context"
	"net/http"
	"time"

	httpkit "github.com/go-kit/kit/transport/http"
	"github.com/google/uuid"
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

// NewCtx will create a new context
func NewCtx() context.Context {
	return context.Background()
}

// CalcTimeoutFromCtx will return the timeout (deadline) for waiting an external response to come back
// TODO: now I return max, need to change it
func CalcTimeoutFromCtx(ctx context.Context) time.Duration {
	return MaxTimeout
}

// GetOrCreateTimeout ...
func GetOrCreateTimeout(ctx context.Context) string {
	val := ctx.Value(TimeoutKey)

	var timeout string
	if val == nil {
		timeout = "15"
	} else {
		timeout = val.(string)
	}

	return timeout
}

// GetOrCreateCorrelationID ...
func GetOrCreateCorrelationID(ctx context.Context) string {
	val := ctx.Value(CorrelationIDKey)

	var corrid string
	if val == nil {
		corrid = uuid.New().String()
	} else {
		corrid = val.(string)
	}

	return corrid
}

// ReadCtx ...
func ReadCtx(ctx context.Context, r *http.Request) context.Context {
	correlationid := r.Header.Get(string(CorrelationIDKey))
	timeout := r.Header.Get(string(TimeoutKey))

	ctx = context.WithValue(ctx, CorrelationIDKey, correlationid)
	ctx = context.WithValue(ctx, TimeoutKey, timeout)

	// need to set the deadline for the context

	return ctx
}

// WriteCtx ...
func WriteCtx(ctx context.Context, r *http.Request) context.Context {
	corrid := GetOrCreateCorrelationID(ctx)
	timeout := GetOrCreateTimeout(ctx)

	r.Header.Add(string(CorrelationIDKey), corrid)
	r.Header.Add(string(timeout), timeout)

	return ctx
}

// WriteCtxBefore ...
func WriteCtxBefore() httpkit.ClientOption {
	return httpkit.ClientBefore(WriteCtx)
}

// ReadCtxBefore ...
func ReadCtxBefore() httpkit.ServerOption {
	return httpkit.ServerBefore(ReadCtx)
}
