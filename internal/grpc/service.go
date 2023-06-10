package grpc

import (
	"context"

	"github.com/barpav/msg-users/internal/data"
)

type Service struct {
	Shutdown chan struct{}
	storage  *data.Storage
}

func (s *Service) Start(storage *data.Storage) error {
	s.storage = storage

	s.Shutdown = make(chan struct{}, 1)

	return nil
}

func (s *Service) Stop(ctx context.Context) (err error) {
	return err
}
