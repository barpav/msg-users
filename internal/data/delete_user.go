package data

import "context"

type queryDeleteUser struct{}

func (q queryDeleteUser) text() string {
	return `
	DELETE FROM users
	WHERE id = $1;
	`
}

func (s *Storage) DeleteUser(ctx context.Context, userId string) (err error) {
	_, err = s.queries[queryDeleteUser{}].ExecContext(ctx, userId)
	return err
}
