package main

import (
	"log"
	"net/http"

	"github.com/barpav/msg-users/internal/rest"
)

func main() {
	service := rest.Service{}

	err := service.Init()

	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(":8080", &service))
}
