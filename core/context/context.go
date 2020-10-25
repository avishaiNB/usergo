package context

import (
	"context"
	"time"
)

// CorrelationIDHeaderKey ..
type CorrelationIDHeaderKey string

// DurationHeaderKey ..
type DurationHeaderKey string

// DeadlineHeaderKey ..
type DeadlineHeaderKey string

const (
	// CorrelationIDKey ...
	CorrelationIDKey CorrelationIDHeaderKey = "correlation_id"

	// DurationKey ..
	DurationKey DurationHeaderKey = "duration"

	// DeadlineKey ...
	DeadlineKey DeadlineHeaderKey = "deadline"
)

// Ctx used to work with context.Context
type Ctx struct {
	Context       context.Context
	Cancel        context.CancelFunc
	CorrelationID string
	Duration      time.Duration
	Deadline      time.Time
}

// NewCtx will create a new context with new corraltion ID, duration and deadline
// Should call it when the service starts
func NewCtx() Ctx {
	c := Ctx{}
	ctx := context.Background()
	corrid := newCorrelation()
	t := NewTimeoutCalculator()
	duration, deadline := t.NewTimeout()
	var cancel context.CancelFunc

	ctx = SetCorrealtionIntoContext(ctx, corrid)
	ctx = SetTimeoutIntoContext(ctx, duration, deadline)
	ctx, cancel = context.WithDeadline(ctx, deadline)

	c.Cancel = cancel
	c.Context = ctx
	c.CorrelationID = corrid
	c.Duration = duration
	c.Deadline = deadline

	return c
}

// SetCorrealtionIntoContext will set into the given context a corraltion ID value
func SetCorrealtionIntoContext(ctx context.Context, correlationID string) context.Context {
	ctx = context.WithValue(ctx, CorrelationIDKey, correlationID)
	return ctx
}

// SetTimeoutIntoContext will set in to the given context the duration and the deadline
func SetTimeoutIntoContext(ctx context.Context, duration time.Duration, deadline time.Time) context.Context {
	ctx = context.WithValue(ctx, DurationKey, duration)
	ctx = context.WithValue(ctx, DeadlineKey, deadline)
	return ctx
}

// GetTimeoutFromContext will return the duration and the deadline from the given context
// If it cannot find it, it will respectively return nil
func GetTimeoutFromContext(ctx context.Context) (time.Duration, time.Time) {
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

// GetOrCreateTimeoutFromContext will get duration and deadline from the context
// if it does not exist it will create new duration and deadline
// If appendToContext, it will update the input context with the duration and deadline
func GetOrCreateTimeoutFromContext(ctx context.Context, appendToContext bool) (time.Duration, time.Time, context.Context) {
	duration, deadline := GetTimeoutFromContext(ctx)

	if deadline.IsZero() {
		t := NewTimeoutCalculator()
		duration, deadline = t.NewTimeout()

		if appendToContext {
			ctx = SetTimeoutIntoContext(ctx, duration, deadline)
		}
	}

	return duration, deadline, ctx
}

// GetOrCreateCorrelationFromContext will get correlation ID from the context
// if it does not exist it will create new correlation ID
// If appendToContext, it will update the input context with the correlation ID
func GetOrCreateCorrelationFromContext(ctx context.Context, appendToContext bool) (string, context.Context) {
	corrid := GetCorrelationID(ctx)

	if corrid == "" {
		corrid = newCorrelation()
		if appendToContext {
			ctx = SetCorrealtionIntoContext(ctx, corrid)
		}
	}

	return corrid, ctx
}

// GetCorrelationID will return the correlation ID from the context
// If it cannot find it, it will return nil
func GetCorrelationID(ctx context.Context) string {
	val := ctx.Value(CorrelationIDKey)

	var corrid string
	if val != nil {
		corrid = val.(string)
	}

	return corrid
}

// GetOrCreateCorrelationID will get the correlation ID from the context
// If it does not exist, it will create a new correlation ID
func GetOrCreateCorrelationID(ctx context.Context) string {
	val := ctx.Value(CorrelationIDKey)

	var corrid string
	if val == nil {
		corrid = newCorrelation()
	} else {
		corrid = val.(string)
	}

	return corrid
}

func newCorrelation() string {
	c := NewCorrelationID()
	return c.New()
}
