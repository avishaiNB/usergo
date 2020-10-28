package amqp

import (
	"context"

	"github.com/streadway/amqp"

	amqptransport "github.com/go-kit/kit/transport/amqp"
	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
	tlerabbitmq "github.com/thelotter-enterprise/usergo/core/transports/rabbitmq"
	"github.com/thelotter-enterprise/usergo/svc/transport"
)

// NewService will create all the rabbitMQ consumers information
// it will not run them.
func NewService(svcEndpoints transport.Endpoints, logger *tlelogger.Manager, connMgr *tlerabbitmq.ConnectionManager) *[]tlerabbitmq.Subscriber {
	subscribers := make([]tlerabbitmq.Subscriber, 0)

	exchangeName := "exchange1"
	queueName := "queue1"
	subscriberName := "command_subscriber"
	loggedInSubscriber := tlerabbitmq.NewCommandSubscriber(
		connMgr,
		subscriberName,
		exchangeName,
		queueName,
		svcEndpoints.UserLoggedInConsumerEndpoint,
		decodeLoggedInUserCommand,
		amqptransport.EncodeJSONResponse,
	)

	// here we can have additional private subscribers
	subscribers = append(subscribers, loggedInSubscriber)
	return &subscribers
}

func decodeLoggedInUserCommand(_ context.Context, d *amqp.Delivery) (interface{}, error) {
	return nil, nil
}
