package transport

import (
	"context"

	tlectx "github.com/thelotter-enterprise/usergo/core/context"
)

// Transport ...
type Transport interface {
	// Read will read from transport into context
	Read(context.Context, interface{}) context.Context

	// Write will write from context into transport
	Write(context.Context, interface{}) context.Context
}

// CreateOutboundContext ...
func CreateOutboundContext(ctx context.Context) (context.Context, context.CancelFunc) {
	calc := tlectx.NewTimeoutCalculator()
	var cancel context.CancelFunc

	_, newCtx := tlectx.GetOrCreateCorrelationFromContext(ctx, true)
	duration, deadline := calc.NextTimeoutFromContext(ctx)

	newCtx = tlectx.SetTimeout(newCtx, duration, deadline)
	newCtx, cancel = context.WithDeadline(newCtx, deadline)

	return newCtx, cancel
}
