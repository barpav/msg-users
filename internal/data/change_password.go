package data

import (
	"context"
)

type queryChangePassword struct{}

func (q queryChangePassword) text() string {
	return `
	UPDATE users SET
		password = MD5($1)::bytea
	WHERE id = $2;
	`
}

func (s *Storage) ChangePassword(ctx context.Context, userId, newPassword string) (err error) {
	_, err = s.queries[queryChangePassword{}].ExecContext(ctx, newPassword, userId)
	return err
}
