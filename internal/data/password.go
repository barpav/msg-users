package data

import (
	"context"
	"crypto/md5"
)

type Password string

const (
	queryPasswordIsValidName = "Password_IsValid"
	queryPasswordIsValid     = `
	SELECT true
	FROM users
	WHERE id = $1 AND password = $2;
	`
)

func (m *Password) IsValid(userId string, s *Storage, ctx context.Context) (result bool, err error) {
	hash := md5.Sum([]byte(*m))

	rows, err := s.queries[queryPasswordIsValidName].QueryContext(ctx, userId, hash[:])
	defer rows.Close()

	return rows.Next(), err
}

// TODO
// func (m *Password) Set(userId string, s *Storage, ctx context.Context) error {
// 	return nil
// }
