package context

import (
	"context"
	"net/http"
	"time"

	httpkit "github.com/go-kit/kit/transport/http"
	"github.com/thelotter-enterprise/usergo/core/utils"
)

// ReadFromHTTPRequest will read from the http request the correlation ID, duration and deadline
// Then it will create a context that reflects the extracted information
// Adding the values and setting the context deadline
func ReadFromHTTPRequest(ctx context.Context, r *http.Request) context.Context {
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

	ctx = SetCorrealtionIntoContext(ctx, correlationID)
	ctx = SetTimeoutIntoContext(ctx, duration, deadline)
	ctx, _ = context.WithDeadline(ctx, deadline)

	return ctx
}

// WriteToHTTPRequest will extract the context values and will append the request as headers
func WriteToHTTPRequest(ctx context.Context, r *http.Request) context.Context {
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
	return httpkit.ClientBefore(WriteToHTTPRequest)
}

// ReadBefore ...
func ReadBefore() httpkit.ServerOption {
	return httpkit.ServerBefore(ReadFromHTTPRequest)
}
