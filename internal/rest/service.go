package rest

import (
	"context"
	"net/http"

	"github.com/rs/zerolog/log"
)

type Service struct {
	Shutdown chan struct{}
	server   *http.Server
	storage  Storage
}

type Storage interface {
	CreateUser(ctx context.Context, id, name, password string) error
}

func (s *Service) Start(storage Storage) {
	s.storage = storage

	s.server = &http.Server{
		Addr:    ":8080",
		Handler: s,
	}

	s.Shutdown = make(chan struct{}, 1)

	go func() {
		err := s.server.ListenAndServe()

		if err != http.ErrServerClosed {
			log.Err(err).Msg("HTTP server crashed.")
		}

		s.Shutdown <- struct{}{}
	}()
}

func (s *Service) Stop(ctx context.Context) (err error) {
	return s.server.Shutdown(ctx)
}

// https://barpav.github.io/msg-api-spec/#/users
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
