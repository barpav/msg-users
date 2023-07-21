package data

import (
	"context"
	"database/sql"

	"github.com/barpav/msg-users/internal/rest/models"
)

type queryGetUserInfo struct{}

func (q queryGetUserInfo) text() string {
	return `
	SELECT name AS name,
	coalesce(picture, '') AS picture 
	FROM users
	WHERE id = $1;
	`
}

func (s *Storage) UserInfo(ctx context.Context, id string) (info *models.UserInfoV1, err error) {
	row := s.queries[queryGetUserInfo{}].QueryRowContext(ctx, id)

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
