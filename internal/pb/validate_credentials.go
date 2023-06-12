package pb

import (
	"context"
	"log"

	"github.com/barpav/msg-users/internal/data"
	usgrpc "github.com/barpav/msg-users/users_service_go_grpc"
)

func (s *Service) Validate(ctx context.Context, credentials *usgrpc.Credentials) (*usgrpc.ValidationResult, error) {

	pass := data.Password(credentials.GetPassword())
	ok, err := pass.IsValid(credentials.GetId(), s.storage, ctx)

	result := &usgrpc.ValidationResult{}

	switch {
	case err != nil:
		result.Status = usgrpc.CredentialsStatus_ERROR
		log.Println(err)
	case ok:
		result.Status = usgrpc.CredentialsStatus_VALID
	default:
		result.Status = usgrpc.CredentialsStatus_NOT_VALID
	}

	return result, err
}
