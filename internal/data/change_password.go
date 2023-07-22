package data

import (
	"context"
	"crypto/md5"
)

type queryChangePassword struct{}

func (q queryChangePassword) text() string {
	return `
	UPDATE users SET
		password = $1
	WHERE id = $2;
	`
}

func (s *Storage) ChangePassword(ctx context.Context, userId, newPassword string) (err error) {
	passwordSum := md5.Sum([]byte(newPassword))
	_, err = s.queries[queryChangePassword{}].ExecContext(ctx, passwordSum[:], userId)
	return err
}
