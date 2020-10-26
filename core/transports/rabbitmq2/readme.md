# RabbitMQ 

### How to create a client? 

In order to create a client, all you need to do is run `Initialize` with the parameters:

`serviceName`: Name of service that is running
`host`: Host name 
`username`: User name
`password`: Password
`logger`: Logger to be used in RabbitMQ


After the rabbitMQ, has been initialize. Now we can access the rabbitmq client using `rabbitmq.Client`

### How to publish ?

In order to publish a message to RabbitMQ, all you need to do is use the `Publish` method and pass the context, the exchangeName and the event you want to publish.

```
rabbitmq.Client.Publish(context.NewContext(), "TheLotter.User.Service" , RandomEvent{})
```

### How to create a consumer ?

In order to define a consumer, all we need to define is a struct with the following methods

```
exchangeName() string
handler(ctx context.Context, message *message) error
```

`ExchangeName() string` will return the name of the exchange you want to listen to

`handler(ctx context.Context, message *message) error` will contain the logic to run for our incomming messages.

> Important: RabbitMQ messages has json format, so when the message will arrive to the consumer, the consumer will need to `json.Unmarshal` in order to pass from json to struct. 

> Keep your struct should define `json:"name"` if there are expected to be unmarshall


Example:

```
type RandomEvent struct {
    name    `json:"name"`
}

type RandomConsumer struct {}

func (r *RandomConsumer) exchangeName() string {
	return "TheLotter.Random.Service:IRandomEvent"
}

func (r *RandomConsumer) handler(ctx context.Context, message *message) error {
    randomEvent := RandomEvent{}
    err := json.Unmarshal(message.Data, &randomEvent)

	// Do something with my random event
}
```



### How to register a consumer ?

Once we have our consumer defined, we need to register it as a private or command consumer.

In order to do that, all we need to do is you the method `RegisterCommandConsumers` or `RegisterPrivateConsumers`.
In this method, we can pass as many consumers as we want.

### How to start consuming ?

After we have defined our consumer and it is register, we just need to run the method `Run()` of the client. 
Now, we should be consuming events from RabbitMQ

