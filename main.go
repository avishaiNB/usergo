package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/thelotter-enterprise/usergo/core"
	"github.com/thelotter-enterprise/usergo/svc"
)

func main() {

	var (
		serviceName string = "user"
		hostAddress string = "localhost:8080"
		zipkinURL   string = "http://localhost:9411/api/v2/spans"
	)

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	errs := make(chan error, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	logger := core.NewLogWithDefaults()
	tracer := core.NewTracer(serviceName, hostAddress, zipkinURL)

	repo := svc.NewRepository()
	service := svc.NewService(logger, tracer, repo)
	endpoints := svc.NewUserEndpoints(logger, tracer, service)
	httpServer := core.NewHTTPServer(logger, tracer, serviceName, hostAddress)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	go func() {
		err := httpServer.Run(&endpoints.HTTPEndpoints)
		if err != nil {
			errs <- err
			fmt.Println(err)
			done <- true
		}
	}()

	<-done
}
