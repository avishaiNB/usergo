package core

import (
	"context"
	"net/http"
	"strconv"
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

// Ctx ..
type Ctx struct {
}

// NewCtx will create a new Ctx
func NewCtx() Ctx {
	return Ctx{}
}

// New ...
func (c Ctx) New() context.Context {
	return context.Background()
}

// CalcTimeout will return the timeout (deadline) for waiting an external response to come back
// TODO: now I return max, need to change it
func (c Ctx) CalcTimeout(ctx context.Context) time.Duration {
	return MaxTimeout
}

// GetOrCreateTimeout ...
func (c Ctx) GetOrCreateTimeout(ctx context.Context) string {
	val := ctx.Value(TimeoutKey)

	var timeout string
	if val == nil {
		sec := c.CalcTimeout(ctx).Seconds()
		timeout = strconv.FormatFloat(sec, 'f', 1, 64)
	} else {
		timeout = val.(string)
	}

	return timeout
}

// GetOrCreateCorrelationID ...
func (c Ctx) GetOrCreateCorrelationID(ctx context.Context) string {
	val := ctx.Value(CorrelationIDKey)

	var corrid string
	if val == nil {
		corrid = NewUUID()
	} else {
		corrid = val.(string)
	}

	return corrid
}

// ReadCtx ...
func (c Ctx) ReadCtx(ctx context.Context, r *http.Request) context.Context {
	correlationid := r.Header.Get(string(CorrelationIDKey))
	timeout := r.Header.Get(string(TimeoutKey))

	ctx = context.WithValue(ctx, CorrelationIDKey, correlationid)
	ctx = context.WithValue(ctx, TimeoutKey, timeout)

	// need to set the deadline for the context

	return ctx
}

// WriteCtx ...
func (c Ctx) WriteCtx(ctx context.Context, r *http.Request) context.Context {
	corrid := c.GetOrCreateCorrelationID(ctx)
	timeout := c.GetOrCreateTimeout(ctx)

	r.Header.Add(string(CorrelationIDKey), corrid)
	r.Header.Add(string(timeout), timeout)

	return ctx
}

// WriteBefore ...
func (c Ctx) WriteBefore() httpkit.ClientOption {
	return httpkit.ClientBefore(c.WriteCtx)
}

// ReadBefore ...
func (c Ctx) ReadBefore() httpkit.ServerOption {
	return httpkit.ServerBefore(c.ReadCtx)
}
