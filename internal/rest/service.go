package rest

import (
	"context"
	"errors"
	"net/http"

	"github.com/barpav/msg-users/internal/data"
)

type Service struct {
	Shutdown chan error
	server   *http.Server
	storage  *data.Storage
}

func (s *Service) Start() error {
	s.storage = &data.Storage{}
	err := s.storage.Open()

	if err != nil {
		return err
	}

	s.server = &http.Server{
		Addr:    ":8080",
		Handler: s,
	}

	s.Shutdown = make(chan error, 1)

	go func() {
		err := s.server.ListenAndServe()

		if !errors.Is(err, http.ErrServerClosed) {
			s.Shutdown <- err
		}
	}()

	return nil
}

func (s *Service) Stop(ctx context.Context) (err error) {
	err = errors.Join(err, s.server.Shutdown(ctx))
	err = errors.Join(err, s.storage.Close(ctx))
	return err
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
