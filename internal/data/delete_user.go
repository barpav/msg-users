package data

import (
	"context"
	"database/sql"
)

func (s *Storage) DeleteUser(ctx context.Context, userId string) (err error) {
	var tx *sql.Tx
	tx, err = s.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "DELETE FROM users WHERE id = $1;", userId)

	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO deleted_users (id) VALUES ($1);", userId)

	if err != nil {
		return err
	}

	return tx.Commit()
}
