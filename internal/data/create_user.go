package data

import (
	"context"
	"crypto/md5"
	"database/sql"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type queryCreateUser struct{}

func (q queryCreateUser) text() string {
	return `
	INSERT INTO users (id, name, password)
	SELECT CAST($1 AS varchar), $2, $3
	WHERE NOT EXISTS (SELECT true FROM deleted_users WHERE id=$1);
	`
}

type ErrUserIdAlreadyExists struct{}

func (s *Storage) CreateUser(ctx context.Context, id, name, password string) (err error) {
	passwordSum := md5.Sum([]byte(password))

	var res sql.Result
	res, err = s.queries[queryCreateUser{}].ExecContext(ctx, id, name, passwordSum[:])

	if err != nil {
		dbErr, ok := err.(*pgconn.PgError)

		if ok && dbErr.Code == pgerrcode.UniqueViolation {
			return &ErrUserIdAlreadyExists{} // existing user id
		}

		return err
	}

	var inserted int64
	inserted, err = res.RowsAffected()

	if err != nil {
		return err
	}

	if inserted == 0 {
		return &ErrUserIdAlreadyExists{} // deleted user id
	}

	return nil
}

func (e *ErrUserIdAlreadyExists) Error() string {
	return "user id already exists"
}

func (e *ErrUserIdAlreadyExists) ImplementsUserIdAlreadyExistsError() {
}
