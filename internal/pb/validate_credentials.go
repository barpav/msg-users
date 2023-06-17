package pb

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/barpav/msg-users/internal/data"
	usgrpc "github.com/barpav/msg-users/users_service_go_grpc"
)

func (s *Service) Validate(ctx context.Context, credentials *usgrpc.Credentials) (*usgrpc.ValidationResult, error) {
	userId := credentials.GetId()

	pass := data.Password(credentials.GetPassword())
	ok, err := pass.IsValid(userId, s.storage, ctx)

	result := &usgrpc.ValidationResult{}

	switch {
	case err != nil:
		result.Status = usgrpc.CredentialsStatus_ERROR
		log.Err(err).Msg(fmt.Sprintf("User '%s' credentials validation failed.", userId))
	case ok:
		result.Status = usgrpc.CredentialsStatus_VALID
		log.Info().Msg(fmt.Sprintf("User '%s' credentials validated (valid).", userId))
	default:
		result.Status = usgrpc.CredentialsStatus_NOT_VALID
		log.Info().Msg(fmt.Sprintf("User '%s' credentials validated (not valid).", userId))
	}

	return result, err
}
