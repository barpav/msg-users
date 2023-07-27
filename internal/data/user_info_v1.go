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
		COALESCE(picture, '') AS picture,
		false AS deleted 
	FROM users
	WHERE id = $1

	UNION ALL

	SELECT
		'' AS name,
		'' AS picture,
		true AS deleted
	FROM deleted_users
	WHERE id = $1;
	`
}

type ErrUserNotFound struct{}
type ErrUserDeleted struct{}

func (s *Storage) UserInfoV1(ctx context.Context, id string) (info *models.UserInfoV1, err error) {
	row := s.queries[queryGetUserInfoV1{}].QueryRowContext(ctx, id)

	info = &models.UserInfoV1{}
	var deleted bool
	err = row.Scan(&info.Name, &info.Picture, &deleted)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &ErrUserNotFound{}
		}

		return nil, err
	}

	if deleted {
		return nil, &ErrUserDeleted{}
	}

	return info, nil
}

func (e *ErrUserNotFound) Error() string {
	return "user not found"
}

func (e *ErrUserNotFound) ImplementsUserNotFoundError() {
}

func (e *ErrUserDeleted) Error() string {
	return "user deleted"
}

func (e *ErrUserDeleted) ImplementsUserDeletedError() {
}
