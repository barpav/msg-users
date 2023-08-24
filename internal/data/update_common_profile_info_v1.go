package data

import (
	"context"

	"github.com/barpav/msg-users/internal/rest/models"
)

type queryUpdateCommonProfileInfoV1 struct{}

func (q queryUpdateCommonProfileInfoV1) text() string {
	return `
	UPDATE users AS current
		SET
		name = COALESCE(NULLIF($1, ''), name),
		picture = COALESCE(NULLIF($2, ''), picture)
	FROM (SELECT picture AS pic FROM users WHERE id = $3) AS previous
	WHERE current.id = $3
	RETURNING
		COALESCE(previous.pic, '') AS old_pic;
	`
}

func (s *Storage) UpdateCommonProfileInfoV1(ctx context.Context, userId string, info *models.UserProfileCommonV1) (oldPic string, err error) {
	row := s.queries[queryUpdateCommonProfileInfoV1{}].QueryRowContext(ctx, info.Name, info.Picture, userId)
	err = row.Scan(&oldPic)
	return oldPic, err
}
