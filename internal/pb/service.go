package pb

import (
	"context"
	"log"
	"net"

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
	s.storage = storage

	s.Shutdown = make(chan struct{}, 1)

	go func() {
		lis, err := net.Listen("tcp", ":9000")

		if err == nil {
			s.server = grpc.NewServer()
			usgrpc.RegisterUsersServer(s.server, s)
			err = s.server.Serve(lis)
		}

		if err != nil {
			log.Println(err)
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
