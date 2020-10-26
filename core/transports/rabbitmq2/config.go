package rabbitmq2

// Config for the rabbitmq client.
type Config struct {
	ConnectionInfo ConnectionInfo

	// Topology responsible of the configuration with rabbitmq when creating,consuming or publishing.
	Topology Topology

	// Marshaller middleware between rabbitmq messages and our events.
	Marshaller Marshaller

	// ErrorHandlers will deal with errors while consuming messages from rabbitmq.
	ErrorHandlers []ErrorHandler

	// Prefetch Count
	PrefetchCount int
}

// NewConfig will return the configuration by default to connect with masstransit
func NewConfig(connectionInfo ConnectionInfo) *Config {
	return &Config{
		ConnectionInfo: connectionInfo,
		Topology:       NewTopology(),
		Marshaller:     &Marshall{},
		ErrorHandlers: []ErrorHandler{
			ErrorQueueHandler{},
		},
		PrefetchCount: 8,
	}
}
