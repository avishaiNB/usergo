package rabbitmq

import (
	"context"
	"time"

	kitamqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
	tlectx "github.com/thelotter-enterprise/usergo/core/context"
	"github.com/thelotter-enterprise/usergo/core/context/transport"
	"github.com/thelotter-enterprise/usergo/core/utils"
)

type amqptransport struct {
}

// NewAMQPTransport will create a new AMQP transport
func NewAMQPTransport() transport.Transport {
	return amqptransport{}
}

func (amqptrans amqptransport) Read(ctx context.Context, req interface{}) (context.Context, context.CancelFunc) {
	pub := req.(*amqp.Publishing)
	corrid := pub.CorrelationId

	ctx = tlectx.SetCorrealtion(ctx, corrid)

	// need to read the deadline and duration and apply deadline
	// ctx.Deadline()

	return ctx, nil
}

func (amqptrans amqptransport) Write(ctx context.Context, req interface{}) (context.Context, context.CancelFunc) {
	newCtx, cancel := transport.CreateTransportContext(ctx)
	corrid := tlectx.GetCorrelation(newCtx)
	duration, deadline := tlectx.GetTimeout(newCtx)

	pub := req.(*amqp.Publishing)
	pub.MessageId = utils.NewUUID()
	pub.Timestamp = utils.NewDateTime().Now()
	pub.CorrelationId = corrid
	pub.ContentType = "application/vnd.masstransit+json"
	headers := setHeaders(pub.Headers, duration, deadline)
	pub.Headers = headers

	return newCtx, cancel
}

// ReadMessageRequestFunc will be executed once a message is consumed
// it will read from the delivery and will create a context
func ReadMessageRequestFunc() kitamqptransport.RequestFunc {
	return func(ctx context.Context, pub *amqp.Publishing, _ *amqp.Delivery) context.Context {
		t := NewAMQPTransport()
		newCtx, _ := t.Read(ctx, pub)
		return newCtx
	}
}

// WriteMessageRequestFunc ...
func WriteMessageRequestFunc() kitamqptransport.RequestFunc {
	return func(ctx context.Context, pub *amqp.Publishing, _ *amqp.Delivery) context.Context {
		t := NewAMQPTransport()
		newCtx, _ := t.Write(ctx, pub)
		return newCtx
	}
}

func setHeaders(headers amqp.Table, duration time.Duration, deadline time.Time) amqp.Table {
	if headers == nil {
		headers = amqp.Table{}
	}

	conv := utils.NewConvertor()
	durationHeader := conv.FromInt64ToString(conv.DurationToMiliseconds(duration))
	deadlineHeader := conv.FromInt64ToString(conv.FromTimeToUnix(deadline))

	headers["tle-deadline-unix"] = deadlineHeader
	headers["tle-duration-ms"] = durationHeader
	headers["tle-caller-process"] = utils.ProcessName()
	headers["tle-caller-hostname"] = utils.HostName()
	headers["tle-caller-processid"] = utils.ProcessID()
	headers["tle-caller-os"] = utils.OperatingSystem()

	return headers
}
