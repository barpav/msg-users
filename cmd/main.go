package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/barpav/msg-users/internal/data"
	"github.com/barpav/msg-users/internal/grpc"
	"github.com/barpav/msg-users/internal/rest"
)

func main() {
	app := microservice{}
	err := app.launch()

	if err == nil {
		log.Println("Microservice is up.")
	} else {
		log.Fatalf("Failed to launch microservice: %s", err)
	}

	err = app.serveAndShutdownGracefully()

	if err == nil {
		log.Println("Microservice stopped.")
	} else {
		log.Printf("Failed to shutdown the microservice gracefully: %s", err)
	}
}

type microservice struct {
	api struct {
		public  *rest.Service // https://barpav.github.io/msg-api-spec/#/users
		private *grpc.Service
	}
	storage *data.Storage
}

func (m *microservice) launch() (err error) {
	m.storage = &data.Storage{}
	err = m.storage.Open()

	if err != nil {
		return err
	}

	m.api.private = &grpc.Service{}
	m.api.private.Start(m.storage)

	m.api.public = &rest.Service{}
	m.api.public.Start(m.storage)

	return err
}

func (m *microservice) serveAndShutdownGracefully() (err error) {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)

	select {
	case <-shutdown:
	case <-m.api.public.Shutdown:
	case <-m.api.private.Shutdown:
	}

	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = errors.Join(err, m.api.public.Stop(ctx))
	err = errors.Join(err, m.api.private.Stop(ctx))
	err = errors.Join(err, m.storage.Close(ctx))

	return err
}
