package data

import (
	"context"
	"crypto/md5"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	queryCreateUserName = "CreateUser"
	queryCreateUser     = `
	INSERT INTO users (id, name, password)
	VALUES ($1, $2, $3);
	`
)

type ErrUserAlreadyExists struct{}

func (s Storage) CreateUser(ctx context.Context, id, name, password string) (err error) {
	passwordSum := md5.Sum([]byte(password))

	_, err = s.queries[queryCreateUserName].ExecContext(ctx, id, name, passwordSum[:])

	if err != nil {
		dbErr, ok := err.(*pgconn.PgError)

		if ok && dbErr.Code == pgerrcode.UniqueViolation {
			return &ErrUserAlreadyExists{}
		}
	}

	return err
}

func (e *ErrUserAlreadyExists) Error() string {
	return "user id already exists"
}

func (e *ErrUserAlreadyExists) ImplementsUserAlreadyExistsError() {
}
