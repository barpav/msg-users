package pb

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	usgrpc "github.com/barpav/msg-users/users_service_go_grpc"
)

func (s *Service) Validate(ctx context.Context, credentials *usgrpc.Credentials) (*usgrpc.ValidationResult, error) {
	userId := credentials.GetId()
	ok, err := s.storage.ValidateCredentials(ctx, userId, credentials.GetPassword())

	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("User '%s' credentials validation failed.", userId))
		return nil, fmt.Errorf("credentials validation failed: %w", err)
	}

	if ok {
		log.Info().Msg(fmt.Sprintf("User '%s' credentials validated (valid).", userId))
	} else {
		log.Info().Msg(fmt.Sprintf("User '%s' credentials validated (not valid).", userId))
	}

	return &usgrpc.ValidationResult{Valid: ok}, nil
}
