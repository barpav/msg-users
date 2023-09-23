package data

import (
	"context"
)

type queryValidateCredentials struct{}

func (q queryValidateCredentials) text() string {
	return `
	SELECT true
	FROM users
	WHERE id = $1 AND password = MD5($2)::bytea;
	`
}

func (s *Storage) ValidateCredentials(ctx context.Context, userId, password string) (valid bool, err error) {
	rows, err := s.queries[queryValidateCredentials{}].QueryContext(ctx, userId, password)

	if err != nil {
		return false, err
	}

	defer rows.Close()

	valid = rows.Next()

	err = rows.Err()

	if err != nil {
		return false, err
	}

	return valid, nil
}
