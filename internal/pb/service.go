package pb

import (
	"context"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"

	usgrpc "github.com/barpav/msg-users/users_service_go_grpc"
	"google.golang.org/grpc"
)

type Service struct {
	Shutdown chan struct{}
	cfg      *Config
	server   *grpc.Server
	storage  Storage

	usgrpc.UnimplementedUsersServer
}

type Storage interface {
	ValidateCredentials(ctx context.Context, userId, password string) (valid bool, err error)
}

func (s *Service) Start(storage Storage) {
	s.cfg = &Config{}
	s.cfg.Read()

	s.server = grpc.NewServer()
	s.storage = storage

	s.Shutdown = make(chan struct{}, 1)

	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.cfg.port))

		if err == nil {
			usgrpc.RegisterUsersServer(s.server, s)
			err = s.server.Serve(lis)
		}

		if err != nil {
			log.Err(err).Msg("gRPC server crashed.")
		}

		s.Shutdown <- struct{}{}
	}()
}

func (s *Service) Stop(ctx context.Context) (err error) {
	closed := make(chan struct{}, 1)

	go func() {
		s.server.GracefulStop()
		closed <- struct{}{}
	}()

	select {
	case <-closed:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
