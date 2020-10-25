package amqp

import (
	"context"

	"github.com/streadway/amqp"

	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
	tlerabbitmq "github.com/thelotter-enterprise/usergo/core/transports/rabbitmq"
	"github.com/thelotter-enterprise/usergo/svc/transport"
)

// NewService ...
func NewService(svcEndpoints transport.Endpoints, logger *tlelogger.Manager) []tlerabbitmq.Consumer {
	consumers := make([]tlerabbitmq.Consumer, 0)

	loggedInConsumer := tlerabbitmq.NewConsumer(
		"user_login_consumer",
		"exchange1",
		"queueq",
		svcEndpoints.UserLoggedInConsumerEndpoint,
		decodeLoggedInUserCommand,
	)

	consumers = append(consumers, loggedInConsumer)
	return consumers
}

func decodeLoggedInUserCommand(_ context.Context, d *amqp.Delivery) (interface{}, error) {
	return nil, nil
}
