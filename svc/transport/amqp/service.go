package amqp

import (
	"context"
	"encoding/json"

	"github.com/streadway/amqp"

	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/thelotter-enterprise/usergo/core/errors"
	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
	tlerabbitmq "github.com/thelotter-enterprise/usergo/core/transports/rabbitmq"
	"github.com/thelotter-enterprise/usergo/core/utils"
	"github.com/thelotter-enterprise/usergo/shared"
	"github.com/thelotter-enterprise/usergo/svc/transport"
)

// NewService will create all the rabbitMQ consumers information
// it will not run them.
func NewService(svcEndpoints transport.Endpoints, logger *tlelogger.Manager, connMgr *tlerabbitmq.ConnectionManager) *[]tlerabbitmq.Subscriber {
	subscribers := make([]tlerabbitmq.Subscriber, 0)

	exchangeName := "exchange1"
	queueName := "queue1"
	subscriberName := "command_subscriber"
	subMgr := tlerabbitmq.NewSubscriberManager(connMgr)

	loggedInSubscriber := subMgr.NewCommandSubscriber(
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

func decodeLoggedInUserCommand(_ context.Context, msg *amqp.Delivery) (interface{}, error) {
	var data shared.LoggedInCommandData
	decoder := utils.NewDecoder()

	m := tlerabbitmq.Message{
		Payload: &tlerabbitmq.MessagePayload{},
	}
	err := json.Unmarshal(msg.Body, &m)
	if err != nil {
		return m, errors.NewApplicationError(err, "failed to decode loggedInUserCommand")
	}
	err = decoder.MapDecode(m.Payload.Data, &data)
	m.Payload.Data = data
	return m, err
}
