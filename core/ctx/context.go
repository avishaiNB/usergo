package ctx

import (
	"context"
	"net/http"
	"time"

	httpkit "github.com/go-kit/kit/transport/http"
	"github.com/thelotter-enterprise/usergo/core/utils"
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

// Ctx used to work with context.Context
type Ctx struct {
	Context       context.Context
	Cancel        context.CancelFunc
	CorrelationID string
	Duration      time.Duration
	Deadline      time.Time
}

// New will create a new context with new corraltion ID, duration and deadline
func New() Ctx {
	ctx := context.Background()
	corrid := NewCorrelation()
	duration, deadline := CalcTimeoutFromContext(ctx)

	return NewFrom(ctx, corrid, duration, deadline)
}

// NewFrom will create a new Ctx
// It will create a new context.Context, add values to the context and set context deadline
func NewFrom(ctx context.Context, correlationID string, duration time.Duration, deadline time.Time) Ctx {
	var cancel context.CancelFunc
	c := Ctx{}
	ctx = SetCorrealtionToContext(ctx, correlationID)
	ctx = SetTimeoutToContext(ctx, duration, deadline)
	ctx, cancel = context.WithDeadline(ctx, deadline)

	c.Cancel = cancel
	c.Context = ctx
	c.CorrelationID = correlationID
	c.Duration = duration
	c.Deadline = deadline

	return c
}

// SetCorrealtionToContext will set into the given context a corraltion ID value
func SetCorrealtionToContext(ctx context.Context, correlationID string) context.Context {
	ctx = context.WithValue(ctx, CorrelationIDKey, correlationID)
	return ctx
}

// SetTimeoutToContext will set in to the given context the duration and the deadline
func SetTimeoutToContext(ctx context.Context, duration time.Duration, deadline time.Time) context.Context {
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

// GetCorrelationFromContext will return the correlation ID from the context
// If it cannot find it, it will return nil
func GetCorrelationFromContext(ctx context.Context) string {
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
func CalcTimeoutFromContext(ctx context.Context) (time.Duration, time.Time) {
	d, t, ctx := GetOrCreateTimeoutFromContext(ctx, false)
	return d, t
}

// NewTimeout will return new timeout including the timeout duration as time.Duration and deadline as time.Time
func NewTimeout() (time.Duration, time.Time) {
	duration := MaxTimeout
	deadline := utils.NewDateTime().AddDuration(duration)

	return duration, deadline
}

// NewCorrelation will return a new correlation ID as string
func NewCorrelation() string {
	return utils.NewUUID()
}

// GetOrCreateTimeoutFromContext ...
func GetOrCreateTimeoutFromContext(ctx context.Context, appendToContext bool) (time.Duration, time.Time, context.Context) {
	duration, deadline := GetTimeoutFromContext(ctx)

	if deadline.IsZero() {
		duration, deadline = NewTimeout()

		if appendToContext {
			ctx = SetTimeoutToContext(ctx, duration, deadline)
		}
	}

	return duration, deadline, ctx
}

// GetOrCreateCorrelationFromContext will get duration and deadline from the context
// if it does not exist it will create new duration and deadline
// If appendToContext, it will update the input context with the duration and deadline
func GetOrCreateCorrelationFromContext(ctx context.Context, appendToContext bool) (string, context.Context) {
	corrid := GetCorrelationFromContext(ctx)

	if corrid == "" {
		corrid = NewCorrelation()
		if appendToContext {
			ctx = SetCorrealtionToContext(ctx, corrid)
		}
	}

	return corrid, ctx
}

// GetOrCreateCorrelationID will get the correlation ID from the context
// If it does not exist, it will create a new correlation ID
// If appendToContext, it will update the input context with the correlation ID
func GetOrCreateCorrelationID(ctx context.Context) string {
	val := ctx.Value(CorrelationIDKey)

	var corrid string
	if val == nil {
		corrid = utils.NewUUID()
	} else {
		corrid = val.(string)
	}

	return corrid
}

// ReadFromRequest will read from the http request the correlation ID, duration and deadline
// Then it will create a context that reflects the extracted information
// Adding the values and setting the context deadline
func ReadFromRequest(ctx context.Context, r *http.Request) context.Context {
	headerCorrelationID := r.Header.Get(string(CorrelationIDKey))
	headerDuration := r.Header.Get(string(DurationKey))
	headerDeadline := r.Header.Get(string(DeadlineKey))

	correlationID := headerCorrelationID
	if headerCorrelationID == "" {
		correlationID = NewCorrelation()
	}

	var duration time.Duration
	var deadline time.Time
	if headerDuration == "" || headerDeadline == "" {
		duration, deadline = NewTimeout()
	} else {
		conv := utils.NewConvertor()
		duration = conv.MilisecondsToDuration(conv.FromStringToInt64(headerDuration))
		deadline = conv.FromUnixToTime(conv.FromStringToInt64(headerDeadline))
	}

	contextFrom := NewFrom(ctx, correlationID, duration, deadline)
	return contextFrom.Context
}

// WriteToRequest will extract the context values and will append the request as headers
func WriteToRequest(ctx context.Context, r *http.Request) context.Context {
	corrid, ctx := GetOrCreateCorrelationFromContext(ctx, false)
	duration, deadline, ctx := GetOrCreateTimeoutFromContext(ctx, false)

	conv := utils.NewConvertor()

	durationHeader := conv.FromInt64ToString(conv.DurationToMiliseconds(duration))
	deadlineHeader := conv.FromInt64ToString(conv.FromTimeToUnix(deadline))
	r.Header.Add(string(CorrelationIDKey), corrid)
	r.Header.Add(string(DurationKey), durationHeader)
	r.Header.Add(string(DeadlineKey), deadlineHeader)

	return ctx
}

// WriteBefore ...
func WriteBefore() httpkit.ClientOption {
	return httpkit.ClientBefore(WriteToRequest)
}

// ReadBefore ...
func ReadBefore() httpkit.ServerOption {
	return httpkit.ServerBefore(ReadFromRequest)
}
