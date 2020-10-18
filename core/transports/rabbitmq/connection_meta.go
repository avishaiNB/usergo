package rabbitmq

import (
	"fmt"
)

// ConnectionMeta ...
type ConnectionMeta struct {
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

// NewConnectionMeta ...
func NewConnectionMeta(host string, port int, username string, password string, vhost string) ConnectionMeta {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", username, password, host, port, vhost)
	return ConnectionMeta{
		URL:         url,
		Host:        host,
		VirtualHost: vhost,
		Pwd:         password,
		Username:    username,
		Port:        port,
	}
}
