package core

import (
	"context"
	"time"

	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
)

// AMQP ..
type AMQP struct {
}

// NewAMQP ..
func NewAMQP() AMQP {
	return AMQP{}
}

// RequestEncoder ..
type RequestEncoder func(context.Context, *amqp.Publishing, interface{}) error

// ResponseDeliveryDecoder ...
type ResponseDeliveryDecoder func(context.Context, *amqp.Delivery) (interface{}, error)

// Publish ...
func (a *AMQP) Publish(ctx context.Context, request interface{}, queueName string) (interface{}, error) {
	corrid := ""
	//var requestEncoder RequestEncoder
	//var responseDeliveryEncoder ResponseDeliveryDecoder
	var channel amqptransport.Channel

	queue := &amqp.Queue{Name: queueName}
	//reqChan := make(chan amqp.Publishing, 1)

	// := &mockChannel{
	// 	f: nullFunc,
	// 	c: reqChan,
	// 	deliveries: []amqp.Delivery{ // we need to mock reply
	// 		amqp.Delivery{
	// 			CorrelationId: corrid,
	// 			//Body:          b,
	// 		},
	// 	},
	// }

	publisher := amqptransport.NewPublisher(
		channel,
		queue,
		func(context.Context, *amqp.Publishing, interface{}) error { return nil },
		func(context.Context, *amqp.Delivery) (response interface{}, err error) {
			return struct{}{}, nil
		},
		// requestEncoder,
		// responseDeliveryEncoder,
		amqptransport.PublisherBefore(amqptransport.SetCorrelationID(corrid)),
		amqptransport.PublisherTimeout(50*time.Second),
	)

	//var publishing amqp.Publishing
	var response interface{}

	var err error
	responseChan := make(chan interface{}, 1)
	errChan := make(chan error, 1)
	go func() {
		res, err := publisher.Endpoint()(ctx, request)
		if err != nil {
			errChan <- err
		} else {
			responseChan <- res
		}
	}()

	select {
	case <-responseChan:
		return response, nil
		break
	case err = <-errChan:
		return response, err
	}

	return response, err
}
