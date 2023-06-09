package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/barpav/msg-users/internal/rest"
)

func main() {
	restService := rest.Service{}
	err := restService.Start()

	if err != nil {
		log.Fatal(err)
	}

	waitAndShutdownGracefully(&restService)
}

func waitAndShutdownGracefully(restService *rest.Service) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)

	select {
	case <-shutdown:
		log.Println("Stopping...")
	case err := <-restService.Shutdown:
		log.Println(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := restService.Stop(ctx)

	if err != nil {
		log.Println(err)
	} else {
		log.Println("Service successfully stopped.")
	}
}
