package http

import (
	"context"
	"net/http"
	"time"

	httpkit "github.com/go-kit/kit/transport/http"
	tlectx "github.com/thelotter-enterprise/usergo/core/context"
	"github.com/thelotter-enterprise/usergo/core/context/manager"
	"github.com/thelotter-enterprise/usergo/core/context/transport"
	"github.com/thelotter-enterprise/usergo/core/utils"
)

type httptransport struct {
}

// NewHTTPTransport will create a new HTTP transport
func NewHTTPTransport() transport.Transport {
	return httptransport{}
}

func (httptrans httptransport) Read(ctx context.Context, req interface{}) context.Context {
	r := req.(*http.Request)

	headerCorrelationID := r.Header.Get(string(tlectx.CorrelationIDKey))
	headerDuration := r.Header.Get(string(tlectx.DurationKey))
	headerDeadline := r.Header.Get(string(tlectx.DeadlineKey))

	correlationID := headerCorrelationID
	if headerCorrelationID == "" {
		correlationID = manager.NewCorrelation()
	}

	var duration time.Duration
	var deadline time.Time
	if headerDuration == "" || headerDeadline == "" {
		t := manager.NewTimeoutCalculator()
		duration, deadline = t.NewTimeout()
	} else {
		conv := utils.NewConvertor()
		duration = conv.MilisecondsToDuration(conv.FromStringToInt64(headerDuration))
		deadline = conv.FromUnixToTime(conv.FromStringToInt64(headerDeadline))
	}

	ctx = manager.SetCorrealtion(ctx, correlationID)
	ctx = manager.SetTimeout(ctx, duration, deadline)
	ctx, _ = context.WithDeadline(ctx, deadline)

	return ctx
}

func (httptrans httptransport) Write(ctx context.Context, req interface{}) context.Context {
	r := req.(*http.Request)
	conv := utils.NewConvertor()

	newCtx, _ := transport.CreateOutboundContext(ctx)
	corrid := manager.GetCorrelation(newCtx)
	duration, deadline := manager.GetTimeout(newCtx)

	durationHeader := conv.FromInt64ToString(conv.DurationToMiliseconds(duration))
	deadlineHeader := conv.FromInt64ToString(conv.FromTimeToUnix(deadline))

	r.Header.Add(string(tlectx.CorrelationIDKey), corrid)
	r.Header.Add(string(tlectx.DurationKey), durationHeader)
	r.Header.Add(string(tlectx.DeadlineKey), deadlineHeader)

	return ctx
}

func write(ctx context.Context, r *http.Request) context.Context {
	t := NewHTTPTransport()
	return t.Write(ctx, r)
}

func read(ctx context.Context, r *http.Request) context.Context {
	t := NewHTTPTransport()
	return t.Read(ctx, r)
}

// WriteBefore ...
func WriteBefore() httpkit.ClientOption {
	return httpkit.ClientBefore(write)
}

// ReadBefore ...
func ReadBefore() httpkit.ServerOption {
	return httpkit.ServerBefore(read)
}
