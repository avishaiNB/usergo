package core

import (
	"context"
	"net/http"
	"time"

	httpkit "github.com/go-kit/kit/transport/http"
)

// CorrelationIDHeaderKey ..
type CorrelationIDHeaderKey string

// DurationHeaderKey ..
type DurationHeaderKey string

// DeadlineHeaderKey ..
type DeadlineHeaderKey string

const (
	// MaxTimeout is 15 seconds
	MaxTimeout time.Duration = time.Second * 15

	// CorrelationIDKey ...
	CorrelationIDKey CorrelationIDHeaderKey = "correlation_id"

	// DurationKey ..
	DurationKey DurationHeaderKey = "duration"

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
	corrid := c.NewCorrelation()
	duration, deadline := c.CalcTimeoutFromContext(ctx)

	return c.NewFrom(ctx, corrid, duration, deadline)
}

// NewFrom ...
func (c *Ctx) NewFrom(ctx context.Context, correlationID string, duration time.Duration, deadline time.Time) Ctx {
	c.SetCorrealtionToContext(ctx, correlationID)
	c.SetTimeoutToContext(ctx, duration, deadline)

	var cancel context.CancelFunc
	ctx, cancel = context.WithDeadline(ctx, deadline)

	c.Cancel = cancel
	c.Context = ctx
	c.CorrelationID = correlationID
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
	ctx = context.WithValue(ctx, DurationKey, duration)
	ctx = context.WithValue(ctx, DeadlineKey, deadline)
	return ctx
}

// GetTimeoutFromContext will return the duration and the deadline from the given context
// If it cannot find it, it will respectively return nil
func (c *Ctx) GetTimeoutFromContext(ctx context.Context) (time.Duration, time.Time) {
	durationAsInterface := ctx.Value(DurationKey)
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

// CalcTimeoutFromContext will return the timeout (deadline) for waiting an external response to come back
// It will return the duration to wait and also the clock timeout
// TODO: now I return max, need to change it
func (c *Ctx) CalcTimeoutFromContext(ctx context.Context) (time.Duration, time.Time) {
	return c.GetOrCreateTimeoutFromContext(ctx, false)
}

// NewTimeout ...
func (c *Ctx) NewTimeout() (time.Duration, time.Time) {
	duration := MaxTimeout
	deadline := NewDateTime().AddDuration(duration)

	return duration, deadline
}

// NewCorrelation ...
func (c *Ctx) NewCorrelation() string {
	return NewUUID()
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

// GetOrCreateCorrelationFromContext will get duration and deadline from the context
// if it does not exist it will create new duration and deadline
// If appendToContext, it will update the input context with the duration and deadline
func (c *Ctx) GetOrCreateCorrelationFromContext(ctx context.Context, appendToContext bool) string {
	corrid := c.GetCorrelationFromContext(ctx)

	if corrid == "" {
		corrid = c.NewCorrelation()
		if appendToContext {
			c.SetCorrealtionToContext(ctx, corrid)
		}
	}

	return corrid
}

// GetOrCreateCorrelationID will get the correlation ID from the context
// If it does not exist, it will create a new correlation ID
// If appendToContext, it will update the input context with the correlation ID
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

// ReadFromRequest will read from the http request the correlation ID, duration and deadline
// Then it will create a context that reflects the extracted information
// Adding the values and setting the context deadline
func (c *Ctx) ReadFromRequest(ctx context.Context, r *http.Request) context.Context {
	headerCorrelationID := r.Header.Get(string(CorrelationIDKey))
	headerDuration := r.Header.Get(string(DurationKey))
	headerDeadline := r.Header.Get(string(DeadlineKey))

	correlationID := headerCorrelationID
	if headerCorrelationID == "" {
		correlationID = c.NewCorrelation()
	}

	var duration time.Duration
	var deadline time.Time
	if headerDuration == "" || headerDeadline == "" {
		duration, deadline = c.NewTimeout()
	} else {
		dt := NewDateTime()
		duration, _ = dt.StringToDuration(headerDuration)
		deadline, _ = dt.StringToTime(headerDeadline)
	}

	contextFrom := c.NewFrom(ctx, correlationID, duration, deadline)
	return contextFrom.Context
}

// WriteToRequest ...
func (c *Ctx) WriteToRequest(ctx context.Context, r *http.Request) context.Context {
	corrid := c.GetOrCreateCorrelationFromContext(ctx, false)
	duration, deadline := c.GetOrCreateTimeoutFromContext(ctx, false)

	dt := NewDateTime()
	durationHeader := dt.DurationToString(duration)
	deadlineHeader := dt.TimeToString(deadline)
	r.Header.Add(string(CorrelationIDKey), corrid)
	r.Header.Add(string(DurationKey), durationHeader)
	r.Header.Add(string(DeadlineKey), deadlineHeader)

	return ctx
}

// WriteBefore ...
func (c *Ctx) WriteBefore() httpkit.ClientOption {
	return httpkit.ClientBefore(c.WriteToRequest)
}

// ReadBefore ...
func (c *Ctx) ReadBefore() httpkit.ServerOption {
	return httpkit.ServerBefore(c.ReadFromRequest)
}
