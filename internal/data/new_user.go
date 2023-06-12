package data

import (
	"context"
	"crypto/md5"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type NewUser struct {
	Id       string
	Name     string
	Password string
}

const (
	queryNewUserCreateName = "NewUser_Create"
	queryNewUserCreate     = `
	INSERT INTO users (id, name, password)
	VALUES ($1, $2, $3);
	`
)

func (m *NewUser) Create(s *Storage, ctx context.Context) (err error, exists bool) {
	passwordHash := md5.Sum([]byte(m.Password))

	_, err = s.queries[queryNewUserCreateName].ExecContext(ctx, m.Id, m.Name, passwordHash[:])

	if err != nil {
		dbErr, ok := err.(*pgconn.PgError)
		exists = ok && dbErr.Code == pgerrcode.UniqueViolation
	}

	return err, exists
}
