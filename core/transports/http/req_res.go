package http

import (
	"context"

	tlecontext "github.com/thelotter-enterprise/usergo/core/context"
	"github.com/thelotter-enterprise/usergo/core/utils"
)

// Response is a base response which will be returned from the transport
type Response struct {
	// Error which occured on the callee
	Error error

	// Data payload returned from the callee
	Data interface{}

	// CircuitOpened is true if execution failed due to circuit being opened
	CircuitOpened bool

	// StatusCode for the execution
	// If using HTTP transport, it will contain HTTP status codes
	StatusCode int
}

// Request is a base response which will be send to the transport
// TODO: use it as the request wrapper
type Request struct {
	// Specific request data
	Data interface{} `mapstructure:",squash"`

	// the correlation ID
	CorrelationID string

	// The duration allowed for the callee to complete the execution
	DurationInMiliseconds int64

	// The deadline for the callee to complete the execution
	DeadlineUnix int64
}

// Wrap will wrap the data in a Request while copying the transport correlation id, duration and timeout
func (r Request) Wrap(ctx context.Context, data interface{}) Request {
	conv := utils.NewConvertor()
	corrid, _ := tlecontext.GetOrCreateCorrelationFromContext(ctx, false)
	// TODO: we need to calculate the deadline and timeout for the callee, so there should be some substruction
	duration, deadline, _ := tlecontext.GetOrCreateTimeoutFromContext(ctx, false)
	req := Request{
		Data:                  data,
		DeadlineUnix:          conv.FromTimeToUnix(deadline),
		DurationInMiliseconds: conv.DurationToMiliseconds(duration),
		CorrelationID:         corrid,
	}
	return req
}
