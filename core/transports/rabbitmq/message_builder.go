package rabbitmq

import (
	"context"
	"fmt"

	tlectx "github.com/thelotter-enterprise/usergo/core/context"
	"github.com/thelotter-enterprise/usergo/core/context/transport"
)

// NewMessage ...
func NewMessage(parentContext context.Context, data interface{}, urn ...string) (Message, context.Context, context.CancelFunc) {
	var urnSlice = make([]string, 0)
	for _, u := range urn {
		urnSlice = append(urnSlice, fmt.Sprintf("urn:message:%v", u))
	}
	payload := MessagePayload{}

	newCtx, cancel := transport.CreateTransportContext(parentContext)
	duration, deadline := tlectx.GetTimeout(newCtx)
	payload.CorrelationID = tlectx.GetCorrelation(newCtx)
	payload.Data = data

	message := Message{
		Payload:       &payload,
		CorrelationID: payload.CorrelationID,
		MessageType:   urnSlice,
		Deadline:      deadline,
		Duration:      duration,
	}

	return message, newCtx, cancel
}
