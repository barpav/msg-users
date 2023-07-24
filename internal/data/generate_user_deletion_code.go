package data

import (
	"context"

	"github.com/google/uuid"
)

type queryGenerateUserDeletionCode struct{}

func (q queryGenerateUserDeletionCode) text() string {
	return `
	INSERT INTO usr_del_confirm_codes (userId, code)
	VALUES ($1, $2)
	ON CONFLICT (userId) DO
		UPDATE SET code = $2;
	`
}

func (s *Storage) GenerateUserDeletionCode(ctx context.Context, userId string) (code string, err error) {
	code = uuid.NewString()
	_, err = s.queries[queryGenerateUserDeletionCode{}].ExecContext(ctx, userId, code)

	if err != nil {
		return "", err
	}

	return code, nil
}
