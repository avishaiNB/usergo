package context

import (
	"context"
	"net/http"
	"time"

	httpkit "github.com/go-kit/kit/transport/http"
	"github.com/thelotter-enterprise/usergo/core/utils"
)

// Transport ...
type Transport interface {

	// ReadFromHTTPRequest will read from the http request the correlation ID, duration and deadline
	// Then it will create a context that reflects the extracted information
	// Adding the values and setting the context deadline
	ReadFromHTTPRequest(context.Context, *http.Request) context.Context

	// WriteToHTTPRequest will extract the context values and will append the request as headers
	WriteToHTTPRequest(context.Context, *http.Request) context.Context

	// CreateOutboundContext should use to create the context which will be used for an outbound call to a service
	// It will create a new context with new corraltion ID, duration and deadline
	CreateOutboundContext(context.Context) (context.Context, context.CancelFunc)
}

type trans struct{}

// NewTransport ...
func NewTransport() Transport {
	return trans{}
}

func (t trans) CreateOutboundContext(ctx context.Context) (context.Context, context.CancelFunc) {
	m := NewManager()
	calc := NewTimeoutCalculator()
	var cancel context.CancelFunc

	_, newCtx := m.GetOrCreateCorrelationFromContext(ctx, true)
	duration, deadline := calc.NextTimeoutFromContext(ctx)

	newCtx = m.SetTimeout(newCtx, duration, deadline)
	newCtx, cancel = context.WithDeadline(newCtx, deadline)

	return newCtx, cancel
}

func (t trans) ReadFromHTTPRequest(ctx context.Context, r *http.Request) context.Context {
	headerCorrelationID := r.Header.Get(string(CorrelationIDKey))
	headerDuration := r.Header.Get(string(DurationKey))
	headerDeadline := r.Header.Get(string(DeadlineKey))

	correlationID := headerCorrelationID
	if headerCorrelationID == "" {
		correlationID = newCorrelation()
	}

	var duration time.Duration
	var deadline time.Time
	if headerDuration == "" || headerDeadline == "" {
		t := NewTimeoutCalculator()
		duration, deadline = t.NewTimeout()
	} else {
		conv := utils.NewConvertor()
		duration = conv.MilisecondsToDuration(conv.FromStringToInt64(headerDuration))
		deadline = conv.FromUnixToTime(conv.FromStringToInt64(headerDeadline))
	}

	m := NewManager()
	ctx = m.SetCorrealtion(ctx, correlationID)
	ctx = m.SetTimeout(ctx, duration, deadline)
	ctx, _ = context.WithDeadline(ctx, deadline)

	return ctx
}

func (t trans) WriteToHTTPRequest(ctx context.Context, r *http.Request) context.Context {
	m := NewManager()
	conv := utils.NewConvertor()

	newCtx, _ := t.CreateOutboundContext(ctx)
	corrid := m.GetCorrelationID(newCtx)
	duration, deadline := m.GetTimeout(newCtx)

	durationHeader := conv.FromInt64ToString(conv.DurationToMiliseconds(duration))
	deadlineHeader := conv.FromInt64ToString(conv.FromTimeToUnix(deadline))

	r.Header.Add(string(CorrelationIDKey), corrid)
	r.Header.Add(string(DurationKey), durationHeader)
	r.Header.Add(string(DeadlineKey), deadlineHeader)

	return ctx
}

// WriteBefore ...
func WriteBefore() httpkit.ClientOption {
	t := NewTransport()
	return httpkit.ClientBefore(t.WriteToHTTPRequest)
}

// ReadBefore ...
func ReadBefore() httpkit.ServerOption {
	t := NewTransport()
	return httpkit.ServerBefore(t.ReadFromHTTPRequest)
}
