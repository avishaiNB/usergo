package manager

import (
	"context"
	"time"

	tlectx "github.com/thelotter-enterprise/usergo/core/context"
)

// Root is the context which is created for a service / application when it initialized
// The context created has no expiration or deadline
// When we need to call external service, we need to create a new context with deadlines, durtion and correlation ID
func Root() context.Context {
	ctx := context.Background()
	corrid := NewCorrelation()
	ctx = SetRootCorrealtion(ctx, corrid)

	return ctx
}

// CreateOutboundContext should use to create the context which will be used for an outbound call to a service
// It will create a new context with new corraltion ID, duration and deadline
func CreateOutboundContext(ctx context.Context) (context.Context, context.CancelFunc) {
	t := NewTimeoutCalculator()
	var cancel context.CancelFunc

	_, newCtx := GetOrCreateCorrelationFromContext(ctx, true)
	duration, deadline := t.NextTimeoutFromContext(ctx)

	newCtx = SetTimeout(newCtx, duration, deadline)
	newCtx, cancel = context.WithDeadline(newCtx, deadline)

	return newCtx, cancel
}

// SetCorrealtion will set into the given context a corraltion ID value
func SetCorrealtion(ctx context.Context, correlationID string) context.Context {
	ctx = context.WithValue(ctx, tlectx.CorrelationIDKey, correlationID)
	return ctx
}

// SetRootCorrealtion will set into the given context a corraltion ID value
func SetRootCorrealtion(ctx context.Context, correlationID string) context.Context {
	ctx = context.WithValue(ctx, tlectx.CorrelationIDRootKey, correlationID)
	return ctx
}

// SetTimeout will set in to the given context the duration and the deadline
func SetTimeout(ctx context.Context, duration time.Duration, deadline time.Time) context.Context {
	ctx = context.WithValue(ctx, tlectx.DurationKey, duration)
	ctx = context.WithValue(ctx, tlectx.DeadlineKey, deadline)
	return ctx
}

// GetTimeout will return the duration and the deadline from the given context
// If it cannot find it, it will respectively return nil
func GetTimeout(ctx context.Context) (time.Duration, time.Time) {
	durationAsInterface := ctx.Value(tlectx.DurationKey)
	deadlineAsInterface := ctx.Value(tlectx.DeadlineKey)

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

// GetCorrelation will return the correlation ID from the context
// if it cannot find it it will try to get the root correlation ID
// If it cannot find it, it will return nil
func GetCorrelation(ctx context.Context) string {
	val := ctx.Value(tlectx.CorrelationIDKey)
	var corrid string

	if val == nil {
		val = GetRootCorrelation(ctx)
	}

	if val != nil {
		corrid = val.(string)
	}

	return corrid
}

// GetRootCorrelation will return the correlation ID from the context
// If it cannot find it, it will return nil
func GetRootCorrelation(ctx context.Context) string {
	val := ctx.Value(tlectx.CorrelationIDRootKey)

	var corrid string
	if val != nil {
		corrid = val.(string)
	}

	return corrid
}

// GetOrCreateCorrelation will get the correlation ID from the context
// If it cannot find, it will try to use the root correlation ID
// If it does not exist, it will create a new correlation ID
func GetOrCreateCorrelation(ctx context.Context) string {
	val := ctx.Value(tlectx.CorrelationIDKey)
	var corrid string

	if val == nil {
		val = GetRootCorrelation(ctx)
		if val == nil {
			corrid = NewCorrelation()
		}
	} else {
		corrid = val.(string)
	}

	return corrid
}

// GetOrCreateTimeout will get duration and deadline from the context
// if it does not exist it will create new duration and deadline
// If appendToContext, it will update the input context with the duration and deadline
func GetOrCreateTimeout(ctx context.Context) (time.Duration, time.Time, context.Context) {
	return GetOrCreateTimeoutFromContext(ctx, false)
}

// GetOrCreateTimeoutFromContext will get duration and deadline from the context
// if it does not exist it will create new duration and deadline
// If appendToContext, it will update the input context with the duration and deadline
func GetOrCreateTimeoutFromContext(ctx context.Context, appendToContext bool) (time.Duration, time.Time, context.Context) {
	duration, deadline := GetTimeout(ctx)

	if deadline.IsZero() {
		t := NewTimeoutCalculator()
		duration, deadline = t.NewTimeout()

		if appendToContext {
			ctx = SetTimeout(ctx, duration, deadline)
		}
	}

	return duration, deadline, ctx
}

// GetOrCreateCorrelationFromContext will get correlation ID from the context
// if it does not exist it will create new correlation ID
// If appendToContext, it will update the input context with the correlation ID
func GetOrCreateCorrelationFromContext(ctx context.Context, appendToContext bool) (string, context.Context) {
	corrid := GetCorrelation(ctx)

	if corrid == "" {
		corrid = NewCorrelation()
		if appendToContext {
			ctx = SetCorrealtion(ctx, corrid)
		}
	}

	return corrid, ctx
}

// NewCorrelation ...
func NewCorrelation() string {
	c := NewCorrelationID()
	return c.New()
}
