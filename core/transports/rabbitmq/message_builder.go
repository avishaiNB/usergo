package rabbitmq

import (
	"context"
	"fmt"

	tlectx "github.com/thelotter-enterprise/usergo/core/context"
	"github.com/thelotter-enterprise/usergo/core/context/transport"
)

// NewMessage will create a rabbit transport message
// It is expected for all messages to be published and consumed from and to rabbit, to be a Message
// The Message includes the payload, context data (such as correlation, timeout)
func NewMessage(parentContext context.Context, data interface{}, urn ...string) (Message, context.Context, context.CancelFunc) {
	var urnSlice = buildURN(urn...)

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

func buildURN(urn ...string) []string {
	var urnSlice = make([]string, 0)
	for _, u := range urn {
		urnSlice = append(urnSlice, fmt.Sprintf("urn:message:%v", u))
	}
	return urnSlice
}
