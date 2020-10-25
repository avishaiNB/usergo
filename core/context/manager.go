package context

import (
	"context"
	"time"
)

// CorrelationIDHeaderKey ..
type CorrelationIDHeaderKey string

// CorrelationIDRootHeaderKey ..
type CorrelationIDRootHeaderKey string

// DurationHeaderKey ..
type DurationHeaderKey string

// DeadlineHeaderKey ..
type DeadlineHeaderKey string

const (
	// CorrelationIDKey ...
	CorrelationIDKey CorrelationIDHeaderKey = "correlation_id"

	// CorrelationIDRootKey ..
	CorrelationIDRootKey CorrelationIDRootHeaderKey = "root_correlation_id"

	// DurationKey ..
	DurationKey DurationHeaderKey = "duration"

	// DeadlineKey ...
	DeadlineKey DeadlineHeaderKey = "deadline"
)

// Manager ..
type Manager interface {
	Root() context.Context
	CreateOutboundContext(context.Context) (context.Context, context.CancelFunc)

	SetTimeout(context.Context, time.Duration, time.Time) context.Context
	GetTimeout(context.Context) (time.Duration, time.Time)
	GetOrCreateTimeoutFromContext(context.Context, bool) (time.Duration, time.Time, context.Context)

	SetCorrealtion(context.Context, string) context.Context
	SetRootCorrealtion(context.Context, string) context.Context
	GetCorrelation(context.Context) string
	GetRootCorrelation(context.Context) string
	GetOrCreateCorrelation(context.Context) string
	GetOrCreateCorrelationFromContext(context.Context, bool) (string, context.Context)
}

// NewManager ...
func NewManager() Manager {
	return ctxmgr{}
}

// Ctx used to work with context.Context
type ctxmgr struct {
}

// Root is the context which is created for a service / application when it initialized
// The context created has no expiration or deadline
// When we need to call external service, we need to create a new context with deadlines, durtion and correlation ID
func (c ctxmgr) Root() context.Context {
	ctx := context.Background()
	corrid := newCorrelation()
	ctx = c.SetRootCorrealtion(ctx, corrid)

	return ctx
}

// CreateOutboundContext should use to create the context which will be used for an outbound call to a service
// It will create a new context with new corraltion ID, duration and deadline
func (c ctxmgr) CreateOutboundContext(ctx context.Context) (context.Context, context.CancelFunc) {
	t := NewTimeoutCalculator()
	var cancel context.CancelFunc

	_, newCtx := c.GetOrCreateCorrelationFromContext(ctx, true)
	duration, deadline := t.NextTimeoutFromContext(ctx)

	newCtx = c.SetTimeout(newCtx, duration, deadline)
	newCtx, cancel = context.WithDeadline(newCtx, deadline)

	return newCtx, cancel
}

// SetCorrealtion will set into the given context a corraltion ID value
func (c ctxmgr) SetCorrealtion(ctx context.Context, correlationID string) context.Context {
	ctx = context.WithValue(ctx, CorrelationIDKey, correlationID)
	return ctx
}

// SetRootCorrealtion will set into the given context a corraltion ID value
func (c ctxmgr) SetRootCorrealtion(ctx context.Context, correlationID string) context.Context {
	ctx = context.WithValue(ctx, CorrelationIDRootKey, correlationID)
	return ctx
}

// SetTimeout will set in to the given context the duration and the deadline
func (c ctxmgr) SetTimeout(ctx context.Context, duration time.Duration, deadline time.Time) context.Context {
	ctx = context.WithValue(ctx, DurationKey, duration)
	ctx = context.WithValue(ctx, DeadlineKey, deadline)
	return ctx
}

// GetTimeout will return the duration and the deadline from the given context
// If it cannot find it, it will respectively return nil
func (c ctxmgr) GetTimeout(ctx context.Context) (time.Duration, time.Time) {
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

// GetCorrelationID will return the correlation ID from the context
// if it cannot find it it will try to get the root correlation ID
// If it cannot find it, it will return nil
func (c ctxmgr) GetCorrelation(ctx context.Context) string {
	val := ctx.Value(CorrelationIDKey)
	var corrid string

	if val == nil {
		val = c.GetRootCorrelation(ctx)
	}

	if val != nil {
		corrid = val.(string)
	}

	return corrid
}

// GetRootCorrelationID will return the correlation ID from the context
// If it cannot find it, it will return nil
func (c ctxmgr) GetRootCorrelation(ctx context.Context) string {
	val := ctx.Value(CorrelationIDRootKey)

	var corrid string
	if val != nil {
		corrid = val.(string)
	}

	return corrid
}

// GetOrCreateCorrelationID will get the correlation ID from the context
// If it cannot find, it will try to use the root correlation ID
// If it does not exist, it will create a new correlation ID
func (c ctxmgr) GetOrCreateCorrelation(ctx context.Context) string {
	val := ctx.Value(CorrelationIDKey)
	var corrid string

	if val == nil {
		val = c.GetRootCorrelation(ctx)
		if val == nil {
			corrid = newCorrelation()
		}
	} else {
		corrid = val.(string)
	}

	return corrid
}

// GetOrCreateTimeoutFromContext will get duration and deadline from the context
// if it does not exist it will create new duration and deadline
// If appendToContext, it will update the input context with the duration and deadline
func (c ctxmgr) GetOrCreateTimeoutFromContext(ctx context.Context, appendToContext bool) (time.Duration, time.Time, context.Context) {
	duration, deadline := c.GetTimeout(ctx)

	if deadline.IsZero() {
		t := NewTimeoutCalculator()
		duration, deadline = t.NewTimeout()

		if appendToContext {
			ctx = c.SetTimeout(ctx, duration, deadline)
		}
	}

	return duration, deadline, ctx
}

// GetOrCreateCorrelationFromContext will get correlation ID from the context
// if it does not exist it will create new correlation ID
// If appendToContext, it will update the input context with the correlation ID
func (c ctxmgr) GetOrCreateCorrelationFromContext(ctx context.Context, appendToContext bool) (string, context.Context) {
	corrid := c.GetCorrelation(ctx)

	if corrid == "" {
		corrid = newCorrelation()
		if appendToContext {
			ctx = c.SetCorrealtion(ctx, corrid)
		}
	}

	return corrid, ctx
}

func newCorrelation() string {
	c := NewCorrelationID()
	return c.New()
}
