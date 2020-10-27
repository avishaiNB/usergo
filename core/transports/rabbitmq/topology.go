package rabbitmq

import "github.com/streadway/amqp"

// Topology responsible of internal configuration of rabbitmq
type Topology interface {
	BuildDurableQueue(channel *amqp.Channel, name string) (amqp.Queue, error)
	BuildNonDurableQueue(ch *amqp.Channel, name string) (amqp.Queue, error)
	BuildDurableExchange(channel *amqp.Channel, name string) error
	BuildNonDurableExchange(channel *amqp.Channel, name string) error
	QueueBind(channel *amqp.Channel, queue, exchange string) error
	Consume(channel *amqp.Channel, queue string) (<-chan amqp.Delivery, error)
	Qos(ch *amqp.Channel) error
	Publish(channel *amqp.Channel, exchange, key string, msg amqp.Publishing) error
}

// topology ...
type topology struct{}

// NewTopology will return the default topology we will use to create rabbitmq object and for publishin and consume
func NewTopology() Topology {
	return topology{}
}

// BuildDurableQueue ...
func (b topology) BuildDurableQueue(ch *amqp.Channel, name string) (amqp.Queue, error) {
	return ch.QueueDeclare(
		name,
		true,  // 	Durable
		false, // 	Auto-delete
		false, //	Exclusive
		false, //	No wait
		nil,   // 	Extra args
	)
}

// BuildNonDurableQueue ...
func (b topology) BuildNonDurableQueue(ch *amqp.Channel, name string) (amqp.Queue, error) {
	return ch.QueueDeclare(
		name,
		false, // 	Durable
		true,  // 	Auto-delete
		true,  //	Exclusive
		false, //	No wait
		nil,   // 	Extra args
	)
}

// BuildDurableExchange ...
func (b topology) BuildDurableExchange(ch *amqp.Channel, name string) error {
	return ch.ExchangeDeclare(
		name,
		"fanout",
		true,  // 	Durable
		false, // 	Auto-delete
		false, // 	Internal
		false, //	No wait
		nil,   // 	Extra args
	)
}

// BuildNonDurableExchange ...
func (b topology) BuildNonDurableExchange(ch *amqp.Channel, name string) error {
	return ch.ExchangeDeclare(
		name,
		"fanout",
		false, // 	Durable
		true,  // 	Auto-delete
		false, // 	Internal
		false, //	No wait
		nil,   // 	Extra args
	)
}

// QueueBind ...
func (b topology) QueueBind(ch *amqp.Channel, queue, exchange string) error {
	return ch.QueueBind(
		queue,
		"", // Key
		exchange,
		false, //	No wait
		nil,   // 	Extra args
	)
}

// Consume ..
func (b topology) Consume(ch *amqp.Channel, queue string) (<-chan amqp.Delivery, error) {
	return ch.Consume(
		queue,
		"",    // 	Consumer Name
		false, // 	Auto Ack
		false, //	Excluvsive
		false, // 	No Local
		false, //	No wait
		nil,   // 	Extra args
	)
}

func (b topology) Qos(ch *amqp.Channel) error {
	return ch.Qos(8, 0, false)
}

// Publish ..
func (b topology) Publish(ch *amqp.Channel, exchange, key string, msg amqp.Publishing) error {
	return ch.Publish(
		exchange, // 	Exchange
		key,      // 	Key
		false,    //	Mandatory
		false,    // 	Inmediate
		msg,      // 	Message
	)
}
