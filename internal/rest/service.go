package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/barpav/msg-users/internal/rest/models"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type Service struct {
	Shutdown chan struct{}
	cfg      *Config
	server   *http.Server
	auth     Authenticator
	storage  Storage
}

type Authenticator interface {
	ValidateSession(ctx context.Context, key, ip, agent string) (userId string, err error)
}

type Storage interface {
	CreateUser(ctx context.Context, id, name, password string) error
	UserInfoV1(ctx context.Context, id string) (*models.UserInfoV1, error)
	UpdateCommonProfileInfoV1(ctx context.Context, userId string, info *models.UserProfileCommonV1) error
	ChangePassword(ctx context.Context, userId, newPassword string) error

	ValidateCredentials(ctx context.Context, userId, password string) (valid bool, err error)
}

func (s *Service) Start(auth Authenticator, storage Storage) {
	s.cfg = &Config{}
	s.cfg.Read()

	s.auth, s.storage = auth, storage

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%s", s.cfg.port),
		Handler: s.operations(),
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

// Specification: https://barpav.github.io/msg-api-spec/#/users
func (s *Service) operations() *chi.Mux {
	ops := chi.NewRouter()

	ops.Use(s.traceInternalServerError)
	ops.Use(s.authenticate)

	// Public endpoint is the concern of the api gateway
	ops.Post("/", s.registerNewUser)
	ops.Get("/", s.getUserInfo)
	ops.Patch("/", s.editUserInfo)
	ops.Delete("/", s.deleteUser)

	return ops
}
