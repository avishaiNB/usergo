// +build integration

package rabbitmq2

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testConsumer struct {
}

func (s testConsumer) exchangeName() string {
	return "1"
}
func (s testConsumer) handler(ctx context.Context, message *Message) error {
	return nil
}

func TestUniqueConsumerForExchange(t *testing.T) {
	username := "user"
	pwd := "pwd"
	host := "localhost"
	vhost := "thelotter"
	port := 5672

	conn := NewConnectionInfo(host, port, username, pwd, vhost)
	config := NewConfig(conn)
	queue := "TheLotter.Skipper.Service-command"

	subscriber, err := newCommandSubscriber(config, queue)
	assert.NoError(t, err)

	err = subscriber.registerConsumer(&testConsumer{}, &testConsumer{})
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errConsumerAlreadyRegisterToExchange))
}
