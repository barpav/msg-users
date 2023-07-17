package data

import (
	"context"
	"crypto/md5"
	"fmt"
)

type queryValidateCredentials struct{}

func (q queryValidateCredentials) text() string {
	return `
	SELECT true
	FROM users
	WHERE id = $1 AND password = $2;
	`
}

func (s *Storage) ValidateCredentials(ctx context.Context, userId, password string) (valid bool, err error) {
	sum := md5.Sum([]byte(password))

	rows, err := s.queries[queryValidateCredentials{}].QueryContext(ctx, userId, sum[:])

	if err != nil {
		return false, fmt.Errorf("Failed to validate credentials: %w", err)
	}

	defer rows.Close()

	return rows.Next(), nil
}
