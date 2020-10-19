package core_test

import (
	"testing"

	"github.com/thelotter-enterprise/usergo/core"
)

func TestNewRabbitMQ(t *testing.T) {
	username := "user"
	pwd := "pwd"
	host := "localhost"
	vhost := "thelotter"
	port := 5672
	log := core.NewLogWithDefaults()
	r := core.NewRabbitMQ(log, host, port, username, pwd, vhost)

	want := "amqp://user:pwd@localhost:5672/thelotter"
	is := r.URL
	if is != want {
		t.Fail()
	}
}
