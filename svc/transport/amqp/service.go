package amqp

import (
	"context"

	"github.com/streadway/amqp"

	amqptransport "github.com/go-kit/kit/transport/amqp"
	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
	"github.com/thelotter-enterprise/usergo/core/transports/rabbitmq"
	tlerabbitmq "github.com/thelotter-enterprise/usergo/core/transports/rabbitmq"
	"github.com/thelotter-enterprise/usergo/svc/transport"
)

// NewService will create all the rabbitMQ consumers information
// it will not run them.
func NewService(svcEndpoints transport.Endpoints, logger *tlelogger.Manager) []tlerabbitmq.Subscriber {
	subscribers := make([]tlerabbitmq.Subscriber, 0)

	options := make([]amqptransport.SubscriberOption, 0)
	exchangeName := "exchange1"
	queueName := "queue1"
	subscriberName := "command_subscriber"
	loggedInSubscriber := tlerabbitmq.NewSubscriber(
		subscriberName,
		exchangeName,
		queueName,
		svcEndpoints.UserLoggedInConsumerEndpoint,
		decodeLoggedInUserCommand,
		amqptransport.EncodeJSONResponse,
		options,
		newUserLoggedInConsumer(exchangeName),
	)

	subscribers = append(subscribers, loggedInSubscriber)
	return subscribers
}

func decodeLoggedInUserCommand(_ context.Context, d *amqp.Delivery) (interface{}, error) {
	return nil, nil
}

type userLoggedInConsumer struct {
}

func newUserLoggedInConsumer(exchangeName string) tlerabbitmq.Consumer {
	return userLoggedInConsumer{}
}

func (c userLoggedInConsumer) MessageURNs() []string {
	return []string{"TheLotter.Service.IUserLoggedInCommand"}
}

func (c userLoggedInConsumer) Name() string {
	return "UserLoggedInConsumer"
}

func (c userLoggedInConsumer) Handler(ctx context.Context, message *rabbitmq.Message) error {
	return nil
}
