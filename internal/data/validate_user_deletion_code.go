package data

import "context"

type queryValidateUserDeletionCode struct{}

func (q queryValidateUserDeletionCode) text() string {
	return `
	SELECT true
	FROM usr_del_confirm_codes
	WHERE userId = $1 AND code = $2;
	`
}

func (s *Storage) ValidateUserDeletionCode(ctx context.Context, userId string, code string) (valid bool, err error) {
	rows, err := s.queries[queryValidateUserDeletionCode{}].QueryContext(ctx, userId, code)

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
