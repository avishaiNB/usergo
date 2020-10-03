package main

import (
	"context"
	"example/user"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
		ctx := context.Background()
		repo := user.NewRepository()
		srv := user.NewService(repo)
		endpoints := user.MakeEndpoints(srv)
		handler := user.NewServer(ctx, endpoints)
		fmt.Println("Listernning on port 8080")
		errs <- http.ListenAndServe(":8080", handler)
	}()

	<-done
}
