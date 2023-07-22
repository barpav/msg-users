package data

import (
	"context"
	"database/sql"

	"github.com/barpav/msg-users/internal/rest/models"
)

type queryGetUserInfoV1 struct{}

func (q queryGetUserInfoV1) text() string {
	return `
	SELECT
		name AS name,
		COALESCE(picture, '') AS picture 
	FROM users
	WHERE id = $1;
	`
}

func (s *Storage) UserInfoV1(ctx context.Context, id string) (info *models.UserInfoV1, err error) {
	row := s.queries[queryGetUserInfoV1{}].QueryRowContext(ctx, id)

	info = &models.UserInfoV1{}
	err = row.Scan(&info.Name, &info.Picture)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return info, nil
}
