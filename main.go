package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/thelotter-enterprise/usergo/core"
	tlamqp "github.com/thelotter-enterprise/usergo/core/amqp"
	tlehttp "github.com/thelotter-enterprise/usergo/core/http"
	"github.com/thelotter-enterprise/usergo/svc"
)

func main() {

	var (
		serviceName      string = "user"
		hostAddress      string = "localhost:8080"
		zipkinURL        string = "http://localhost:9411/api/v2/spans"
		rabbitMQUsername string = "thelotter"
		rabbitMQPwd      string = "Dhvbuo1"
		rabbitMQHost     string = "int-k8s1"
		rabbitMQVhost    string = "thelotter"
		rabbitMQPort     int    = 32672
	)

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	errs := make(chan error, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	logger := core.NewLogWithDefaults()
	tracer := core.NewTracer(serviceName, hostAddress, zipkinURL)
	rabbitmq := tlamqp.NewRabbitMQ(logger, rabbitMQHost, rabbitMQPort, rabbitMQUsername, rabbitMQPwd, rabbitMQVhost)

	repo := svc.NewRepository()
	service := svc.NewService(logger, tracer, repo)
	httpEndpoints := svc.NewUserHTTPEndpoints(logger, tracer, service)
	httpServer := tlehttp.NewHTTPServer(logger, tracer, serviceName, hostAddress)

	amqpEndpoints := svc.NewUserAMQPConsumerEndpoints(logger, tracer, service, &rabbitmq)
	amqpServer := tlamqp.NewAMQPServer(logger, tracer, &rabbitmq, serviceName)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	go func() {
		err := httpServer.Run(httpEndpoints.HTTPEndpoints)
		if err != nil {
			errs <- err
			fmt.Println(err)
			done <- true
		}
	}()

	go func() {
		err := amqpServer.Run(amqpEndpoints.Consumers)
		if err != nil {
			errs <- err
			fmt.Println(err)
			done <- true
		}
	}()

	<-done
}
