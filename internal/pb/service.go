package pb

import (
	"context"
	"net"

	"github.com/rs/zerolog/log"

	"github.com/barpav/msg-users/internal/data"
	usgrpc "github.com/barpav/msg-users/users_service_go_grpc"
	"google.golang.org/grpc"
)

type Service struct {
	Shutdown chan struct{}
	server   *grpc.Server
	storage  *data.Storage

	usgrpc.UnimplementedUsersServer
}

func (s *Service) Start(storage *data.Storage) {
	s.server = grpc.NewServer()
	s.storage = storage

	s.Shutdown = make(chan struct{}, 1)

	go func() {
		lis, err := net.Listen("tcp", ":9000")

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
