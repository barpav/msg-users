package data

import (
	"context"

	"github.com/barpav/msg-users/internal/rest/models"
)

type queryUpdateCommonProfileInfoV1 struct{}

func (q queryUpdateCommonProfileInfoV1) text() string {
	return `
	UPDATE users SET
		name = COALESCE(NULLIF($1, ''), name),
		picture = COALESCE(NULLIF($2, ''), picture)
	WHERE id = $3;
	`
}

func (s *Storage) UpdateCommonProfileInfoV1(ctx context.Context, userId string, info *models.UserProfileCommonV1) (err error) {
	_, err = s.queries[queryUpdateCommonProfileInfoV1{}].ExecContext(ctx, info.Name, info.Picture, userId)
	return err
}
