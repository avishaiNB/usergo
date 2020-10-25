package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
	tlemetrics "github.com/thelotter-enterprise/usergo/core/metrics"
	tletracer "github.com/thelotter-enterprise/usergo/core/tracer"
	tlerabbitmq "github.com/thelotter-enterprise/usergo/core/transports/rabbitmq"
	tleimpl "github.com/thelotter-enterprise/usergo/svc/implementation"
	svcmw "github.com/thelotter-enterprise/usergo/svc/middleware"
	svctrans "github.com/thelotter-enterprise/usergo/svc/transport"
	svcamqp "github.com/thelotter-enterprise/usergo/svc/transport/amqp"
	svchttp "github.com/thelotter-enterprise/usergo/svc/transport/http"
)

func main() {
	var (
		serviceName      string                    = "user"
		hostAddress      string                    = "localhost:8080"
		zipkinURL        string                    = "http://localhost:9411/api/v2/spans"
		rabbitMQUsername string                    = "thelotter"
		rabbitMQPwd      string                    = "Dhvbuo1"
		rabbitMQHost     string                    = "int-k8s1"
		rabbitMQVhost    string                    = "thelotter"
		rabbitMQPort     int                       = 32672
		env              string                    = "dev"
		logLevel         tlelogger.AtomicLevelName = tlelogger.Debug
	)

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	errs := make(chan error, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// Setting up the infra services which will be used
	stdConf := tlelogger.Config{
		LevelName:  logLevel,
		Env:        env,
		LoggerName: "std",
	}
	stdLogger := tlelogger.NewStdOutLogger(stdConf)
	logManager := tlelogger.NewLoggerManager(stdLogger)

	tracer := tletracer.NewTracer(serviceName, hostAddress, zipkinURL)
	inst := tlemetrics.NewPrometheusInstrumentor(serviceName)

	// implementation
	repo := tleimpl.NewRepository()
	service := tleimpl.NewService(&logManager, tracer, repo)

	// middlewares
	service = svcmw.NewLoggingMiddleware(&logManager)(service)             // Hook up the logging middleware
	service = svcmw.NewInstrumentingMiddleware(&logManager, inst)(service) // Hook up the inst middleware

	// http
	httpEndpoints := svctrans.MakeEndpoints(service)
	handler := svchttp.NewService(httpEndpoints, make([]kithttp.ServerOption, 0), logManager.(log.Logger))
	go func() {
		server := &http.Server{
			Addr:    hostAddress,
			Handler: handler,
		}
		errs <- server.ListenAndServe()
		done <- true
	}()

	// setting up RabbitMQ server
	conn := tlerabbitmq.NewConnectionInfo(rabbitMQHost, rabbitMQPort, rabbitMQUsername, rabbitMQPwd, rabbitMQVhost)
	rabbitmq := tlerabbitmq.NewRabbitMQ(&logManager, conn)
	amqpEndpoints := svcamqp.NewUserAMQPConsumerEndpoints(&logManager, tracer, service, &rabbitmq)
	amqpServer := tlerabbitmq.NewServer(&logManager, tracer, &rabbitmq, serviceName)

	go func() {
		err := amqpServer.Run(amqpEndpoints.Consumers)
		if err != nil {
			errs <- err
			fmt.Println(err)
			done <- true
		}
	}()

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	<-done
}
