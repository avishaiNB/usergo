package amqp

import (
	"fmt"

	"github.com/streadway/amqp"
)

// Connection ...
type Connection struct {
	// URL like amqp://guest:guest@localhost:5672/
	URL string

	// Usewrname to connect to RabbitMQ
	Username string

	// Pwd to connect to RabbitMQ
	Pwd string

	// VirtualHost to connect to RabbitMQ
	VirtualHost string

	// Port to connect to RabbitMQ
	Port int

	// Host to connect to RabbitMQ
	Host string
}

// NewConnection ...
func NewConnection(host string, port int, username string, password string, vhost string) Connection {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", username, password, host, port, vhost)
	return Connection{
		URL:         url,
		Host:        host,
		VirtualHost: vhost,
		Pwd:         password,
		Username:    username,
		Port:        port,
	}
}

// Connect will create a new connection to RabbitMQ based on the input entered when created the RabbitMQ instance
// Connection will be returned BUT also stored in the RabbitMQ instance
func (a *RabbitMQ) Connect() (*amqp.Connection, error) {
	if a.AMQPConnection != nil {
		return a.AMQPConnection, nil
	}
	conn, err := amqp.Dial(a.Connection.URL)
	if err == nil {
		a.AMQPConnection = conn
	}
	return conn, err
}

// Close will close the open connection attached to the RabbitMQ instance
func (a *RabbitMQ) Close() error {
	var err error
	if a.AMQPConnection != nil {
		err = a.AMQPConnection.Close()
	}
	return err
}
