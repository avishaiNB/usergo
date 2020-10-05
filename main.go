package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/thelotter-enterprise/usergo/svc"
)

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	errs := make(chan error)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	go func() {
		repo := svc.NewRepository()
		srv := svc.NewService(repo)
		endpoints := svc.MakeEndpoints(srv)
		server := svc.MakeServer("user", "localhost:8080", "http://localhost:9411/api/v2/spans", endpoints, errs)
		server.Run()
	}()

	<-done
}
