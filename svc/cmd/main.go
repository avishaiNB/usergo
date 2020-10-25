package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kithttp "github.com/go-kit/kit/transport/http"
	tlectx "github.com/thelotter-enterprise/usergo/core/context"
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
		ctx              context.Context           = tlectx.NewCtxManager().Root()
	)

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	errs := make(chan error, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	// infra services
	stdConf := tlelogger.Config{LevelName: logLevel, Env: env, LoggerName: "std"}
	stdLogger := tlelogger.NewStdOutLogger(stdConf)
	logManager := tlelogger.NewLoggerManager(stdLogger)
	tracer := tletracer.NewTracer(serviceName, hostAddress, zipkinURL)
	promMetricsInst := tlemetrics.NewPrometheusInstrumentor(serviceName)

	// implementation
	repo := tleimpl.NewRepository()
	service := tleimpl.NewService(&logManager, tracer, repo)

	// middlewares
	service = svcmw.NewLoggingMiddleware(&logManager)(service)                        // Hook up the logging middleware
	service = svcmw.NewInstrumentingMiddleware(&logManager, promMetricsInst)(service) // Hook up the inst middleware

	// making all types of endpoints
	endpoints := svctrans.MakeEndpoints(service)

	// http
	handler := svchttp.NewService(ctx, endpoints, make([]kithttp.ServerOption, 0), logManager)
	go func() {
		server := &http.Server{
			Addr:    hostAddress,
			Handler: handler,
		}
		logManager.Info(ctx, fmt.Sprintf("listening for http calls on %s", hostAddress))
		errs <- server.ListenAndServe()
		done <- true
	}()

	// setting up RabbitMQ server
	conn := tlerabbitmq.NewConnectionInfo(rabbitMQHost, rabbitMQPort, rabbitMQUsername, rabbitMQPwd, rabbitMQVhost)
	rabbitmq := tlerabbitmq.NewRabbitMQ(&logManager, conn)
	consumers := svcamqp.NewService(endpoints, &logManager)
	amqpServer := tlerabbitmq.NewServer(&logManager, tracer, &rabbitmq)

	go func() {
		logManager.Info(ctx, fmt.Sprintf("listening for amqp messages"))
		err := amqpServer.Run(ctx, &consumers)
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
