package data

import (
	"context"
	"crypto/md5"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Storage) prepareRegisterQuery() (err error) {
	const query = `
	INSERT INTO users (id, name, password)
	VALUES ($1, $2, $3);
	`
	s.queries.register, err = s.db.Prepare(query)

	return err
}

type NewUser struct {
	Id       string
	Name     string
	Password string
}

func (s *Storage) Register(user *NewUser, ctx context.Context) (err error, exists bool) {
	passwordHash := md5.Sum([]byte(user.Password))

	_, err = s.queries.register.ExecContext(ctx, user.Id, user.Name, passwordHash[:])

	if err != nil {
		dbErr, ok := err.(*pgconn.PgError)
		exists = ok && dbErr.Code == pgerrcode.UniqueViolation
	}

	return err, exists
}
