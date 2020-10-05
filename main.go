package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
	errs := make(chan error)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	logger := svc.NewLogger()
	tracer := svc.NewTracer(serviceName, hostAddress, zipkinURL)

	repo := svc.NewRepository()
	service := svc.NewService(&logger, &tracer, repo)
	endpoints := svc.NewEndpoints(&logger, &tracer, &service)
	server := svc.NewServer(&logger, &tracer, serviceName, hostAddress, errs)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	go func() {
		server.Run(&endpoints)
	}()

	<-done
}
