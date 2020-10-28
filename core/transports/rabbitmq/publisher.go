package rabbitmq

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-kit/kit/endpoint"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
	tlectx "github.com/thelotter-enterprise/usergo/core/context"
	"github.com/thelotter-enterprise/usergo/core/errors"
	"github.com/thelotter-enterprise/usergo/core/utils"
)

// Publisher ...
type Publisher interface {
	PublishEndpoint(context.Context, string, amqptransport.EncodeRequestFunc, ...amqptransport.PublisherOption) (endpoint.Endpoint, error)
	Close(context.Context) error
}

type publisher struct {
	connectionManager *ConnectionManager
	ch                *amqp.Channel
	isConnected       bool
	topology          Topology
}

// NewPublisher will create a new publisher and will establish a connection to rabbit
func NewPublisher(conn *ConnectionManager) Publisher {
	p := publisher{
		connectionManager: conn,
		topology:          NewTopology(),
	}
	p.connect()
	return &p
}

func (p publisher) PublishEndpoint(ctx context.Context, exchangeName string, encodeFunc amqptransport.EncodeRequestFunc, options ...amqptransport.PublisherOption) (endpoint.Endpoint, error) {
	if p.isConnected == false {
		return nil, errors.NewApplicationErrorf("before publishing, you must connect to rabbitMQ")
	}

	p.buildExchange(exchangeName)

	corrid := tlectx.GetCorrelation(ctx)
	duration, deadline := tlectx.GetTimeout(ctx)

	var queue *amqp.Queue
	queue = &amqp.Queue{Name: ""}

	// building the publisher options
	ops := make([]amqptransport.PublisherOption, 0)
	ops = append(ops, options...)
	before := amqptransport.PublisherBefore(
		setMessageID(),
		setMessageTimestamp(),
		setHeaders(deadline, duration),
		amqptransport.SetContentType("application/vnd.masstransit+json"),
		amqptransport.SetCorrelationID(corrid),
		amqptransport.SetPublishDeliveryMode(2),
		amqptransport.SetPublishExchange(exchangeName),
	)
	ops = append(ops, before)
	ops = append(ops,
		amqptransport.PublisherTimeout(duration),
		amqptransport.PublisherDeliverer(amqptransport.SendAndForgetDeliverer),
	)

	publisher := amqptransport.NewPublisher(p.ch, queue, encodeFunc, NoopResponseDecoder, ops...)

	return publisher.Endpoint(), nil
}

func (p publisher) buildExchange(exchanegName string) {
	conn := *p.connectionManager
	ch, err := conn.GetChannel()
	if err == nil {
		defer ch.Close()
		p.topology.BuildDurableExchange(ch, exchanegName)
	}
}

// NoopResponseDecoder is a no operation needed
// Used for One way messages
func NoopResponseDecoder(ctx context.Context, d *amqp.Delivery) (response interface{}, err error) {
	return struct{}{}, nil
}

// DefaultRequestEncoder ...
func DefaultRequestEncoder(exchangeName string) func(context.Context, *amqp.Publishing, interface{}) error {
	f := func(ctx context.Context, p *amqp.Publishing, request interface{}) error {
		message := request.(*Message)
		body, err := json.Marshal(message)
		p.Body = body
		return err
	}
	return f
}

// Close will shutdown the client gracely
func (p *publisher) Close(ctx context.Context) error {
	var err error

	if p.isConnected && p.ch != nil {
		cherr := p.ch.Close()

		if cherr != nil {
			err = errors.NewApplicationError(err, cherr.Error())
		} else {
			p.isConnected = false
		}
	}

	return err
}

func (p *publisher) connect() error {
	connMgr := *p.connectionManager
	ch, err := connMgr.GetChannel()
	if err == nil {
		p.ch = ch
		//p.changeConnection(ctx, conn, ch)
		p.isConnected = true
	}
	return err
}

func setMessageID() amqptransport.RequestFunc {
	return func(ctx context.Context, pub *amqp.Publishing, _ *amqp.Delivery) context.Context {
		pub.MessageId = utils.NewUUID()
		return ctx
	}
}

func setMessageTimestamp() amqptransport.RequestFunc {
	return func(ctx context.Context, pub *amqp.Publishing, _ *amqp.Delivery) context.Context {
		pub.Timestamp = utils.NewDateTime().Now()
		return ctx
	}
}

func setHeaders(deadline time.Time, duration time.Duration) amqptransport.RequestFunc {
	return func(ctx context.Context, pub *amqp.Publishing, _ *amqp.Delivery) context.Context {
		headers := pub.Headers
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

		pub.Headers = headers
		return ctx
	}
}
