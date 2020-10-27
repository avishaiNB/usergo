package rabbitmq

import (
	"github.com/streadway/amqp"
	tleerrors "github.com/thelotter-enterprise/usergo/core/errors"
)

// NewChannel will create a new rabbitMQ channel
func (a *Client) NewChannel() (*amqp.Channel, error) {
	var err error
	var ch *amqp.Channel
	if a.AMQPConnection == nil {
		err = tleerrors.New("Connect to rabbit before tring to get a channel")
		// TODO: better logging here
		//a.LogManager.Error()
	} else {
		ch, err = a.AMQPConnection.Channel()
		if err != nil {
			// TODO: better logging here
			//a.LogManager.Logger.Log(err)
		}
	}

	return ch, err
}
