package svc

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tlelogger "github.com/thelotter-enterprise/usergo/core/logger"
	tlemetrics "github.com/thelotter-enterprise/usergo/core/metrics"
	tletracer "github.com/thelotter-enterprise/usergo/core/tracer"
	tlehttp "github.com/thelotter-enterprise/usergo/core/transports/http"
	tlerabbitmq "github.com/thelotter-enterprise/usergo/core/transports/rabbitmq"
)

// Run ...
func Run() {
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

	// In this part we are building the service and extending it using middleware pattern
	repo := NewRepository()
	service := NewService(&logManager, tracer, repo)
	service = NewLoggingMiddleware(&logManager)(service)             // Hook up the logging middleware
	service = NewInstrumentingMiddleware(&logManager, inst)(service) // Hook up the inst middleware

	// setting up the http server
	httpEndpoints := NewUserHTTPEndpoints(&logManager, tracer, service)
	httpServer := tlehttp.NewServer(&logManager, tracer, serviceName, hostAddress)

	// setting up RabbitMQ server
	conn := tlerabbitmq.NewConnectionMeta(rabbitMQHost, rabbitMQPort, rabbitMQUsername, rabbitMQPwd, rabbitMQVhost)
	rabbitmq := tlerabbitmq.NewRabbitMQ(&logManager, conn)
	amqpEndpoints := NewUserAMQPConsumerEndpoints(&logManager, tracer, service, &rabbitmq)
	amqpServer := tlerabbitmq.NewServer(&logManager, tracer, &rabbitmq, serviceName)

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
