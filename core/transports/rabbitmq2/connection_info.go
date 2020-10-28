package rabbitmq2

import (
	"fmt"
)

// ConnectionInfo ...
type ConnectionInfo struct {
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

// NewConnectionInfo ...
func NewConnectionInfo(host string, port int, username string, password string, vhost string) ConnectionInfo {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", username, password, host, port, vhost)
	return ConnectionInfo{
		URL:         url,
		Host:        host,
		VirtualHost: vhost,
		Pwd:         password,
		Username:    username,
		Port:        port,
	}
}
