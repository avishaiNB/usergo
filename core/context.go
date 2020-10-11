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

// DeadlineHeaderKey ..
type DeadlineHeaderKey string

const (
	// MaxTimeout is 15 seconds
	MaxTimeout time.Duration = time.Second * 15

	// CorrelationIDKey ...
	CorrelationIDKey CorrelationIDHeaderKey = "correlation_id"

	// TimeoutKey ..
	TimeoutKey TimeoutHeaderKey = "timeout"

	// DeadlineKey ...
	DeadlineKey DeadlineHeaderKey = "deadline"
)

// Ctx ..
type Ctx struct {
	Context       context.Context
	Cancel        context.CancelFunc
	CorrelationID string
	Duration      time.Duration
	Deadline      time.Time
}

// NewCtx will create a new Ctx
func NewCtx() Ctx {
	return Ctx{}
}

// New will create a new context with new corraltion ID, duration and deadline
func (c *Ctx) New() Ctx {
	ctx := context.Background()
	corrid := NewUUID()
	duration, deadline := c.CalcTimeout(ctx)

	c.SetCorrealtionToContext(ctx, corrid)
	c.SetTimeoutToContext(ctx, duration, deadline)

	var cancel context.CancelFunc
	ctx, cancel = context.WithDeadline(ctx, deadline)

	c.Cancel = cancel
	c.Context = ctx
	c.CorrelationID = corrid
	c.Duration = duration
	c.Deadline = deadline

	return *c
}

// SetCorrealtionToContext will set into the given context a corraltion ID value
func (c *Ctx) SetCorrealtionToContext(ctx context.Context, correlationID string) context.Context {
	ctx = context.WithValue(ctx, CorrelationIDKey, correlationID)
	return ctx
}

// SetTimeoutToContext will set in to the given context the duration and the deadline
func (c *Ctx) SetTimeoutToContext(ctx context.Context, duration time.Duration, deadline time.Time) context.Context {
	ctx = context.WithValue(ctx, TimeoutKey, duration)
	ctx = context.WithValue(ctx, DeadlineKey, deadline)
	return ctx
}

// GetTimeoutFromContext will return the duration and the deadline from the given context
// If it cannot find it, it will respectively return nil
func (c *Ctx) GetTimeoutFromContext(ctx context.Context) (time.Duration, time.Time) {
	durationAsInterface := ctx.Value(TimeoutKey)
	deadlineAsInterface := ctx.Value(DeadlineKey)

	var duration time.Duration
	var deadline time.Time

	if durationAsInterface != nil {
		duration = durationAsInterface.(time.Duration)
	}
	if deadlineAsInterface != nil {
		deadline = deadlineAsInterface.(time.Time)
	}

	return duration, deadline
}

// GetCorrelationFromContext will return the correlation ID from the context
// If it cannot find it, it will return nil
func (c *Ctx) GetCorrelationFromContext(ctx context.Context) string {
	val := ctx.Value(CorrelationIDKey)

	var corrid string
	if val != nil {
		corrid = val.(string)
	}

	return corrid
}

// CalcTimeout will return the timeout (deadline) for waiting an external response to come back
// It will return the duration to wait and also the clock timeout
// TODO: now I return max, need to change it
func (c *Ctx) CalcTimeout(ctx context.Context) (time.Duration, time.Time) {
	duration := MaxTimeout
	deadline := NewDateTime().AddDuration(duration)

	return duration, deadline
}

// NewTimeout ...
func (c *Ctx) NewTimeout() (time.Duration, time.Time) {
	duration := MaxTimeout
	deadline := NewDateTime().AddDuration(duration)

	return duration, deadline
}

// GetOrCreateTimeoutFromContext ...
func (c *Ctx) GetOrCreateTimeoutFromContext(ctx context.Context, appendToContext bool) (time.Duration, time.Time) {
	duration, deadline := c.GetTimeoutFromContext(ctx)

	if deadline.IsZero() {
		duration, deadline = c.NewTimeout()

		if appendToContext {
			c.SetTimeoutToContext(ctx, duration, deadline)
		}
	}

	return duration, deadline
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
