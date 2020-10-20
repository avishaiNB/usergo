package rabbitmq

import (
	"github.com/streadway/amqp"
	tleerrors "github.com/thelotter-enterprise/usergo/core/errors"
)

// NewChannel will create a new rabbitMQ channel
func (a *RabbitMQ) NewChannel() (*amqp.Channel, error) {
	var err error
	var ch *amqp.Channel
	if a.AMQPConnection == nil {
		err = tleerrors.NewApplicationError("Connect to rabbit before tring to get a channel")
		// TODO: better logging here
		a.Log.Logger.Log(err)
	} else {
		ch, err = a.AMQPConnection.Channel()
		if err != nil {
			// TODO: better logging here
			a.Log.Logger.Log(err)
		}
	}

	return ch, err
}
