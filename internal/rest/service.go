package rest

import (
	"context"
	"log"
	"net/http"

	"github.com/barpav/msg-users/internal/data"
)

type Service struct {
	Shutdown chan struct{}
	server   *http.Server
	storage  *data.Storage
}

func (s *Service) Start(storage *data.Storage) {
	s.storage = storage

	s.server = &http.Server{
		Addr:    ":8080",
		Handler: s,
	}

	s.Shutdown = make(chan struct{}, 1)

	go func() {
		log.Println(s.server.ListenAndServe())
		s.Shutdown <- struct{}{}
	}()
}

func (s *Service) Stop(ctx context.Context) (err error) {
	return s.server.Shutdown(ctx)
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// URL handling is reverse proxy's concern
	switch r.Method {
	case http.MethodPost:
		s.registerNewUser(w, r)
	case http.MethodGet:
		s.getUserInfo(w, r)
	case http.MethodPut:
		s.editUserInfo(w, r)
	case http.MethodDelete:
		s.deleteUser(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
