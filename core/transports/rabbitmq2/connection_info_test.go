package rabbitmq2_test

import (
	"testing"

	rabbitmq "github.com/thelotter-enterprise/usergo/core/transports/rabbitmq2"
)

func TestConnectionInfoURL(t *testing.T) {
	username := "user"
	pwd := "pwd"
	host := "localhost"
	vhost := "thelotter"
	port := 5672

	connectionMeta := rabbitmq.NewConnectionInfo(host, port, username, pwd, vhost)

	want := "amqp://user:pwd@localhost:5672/thelotter"
	is := connectionMeta.URL
	if is != want {
		t.Fail()
	}
}
