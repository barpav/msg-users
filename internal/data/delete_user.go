package data

import (
	"context"
	"database/sql"
)

type queryDeleteUser struct{}

func (q queryDeleteUser) text() string {
	return `
	DELETE FROM users WHERE id = $1;
	`
}

type queryReserveDeletedUserId struct{}

func (q queryReserveDeletedUserId) text() string {
	return `
	INSERT INTO deleted_users (id) VALUES ($1);
	`
}

func (s *Storage) DeleteUser(ctx context.Context, userId string) (err error) {
	var tx *sql.Tx
	tx, err = s.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	defer tx.Rollback()

	_, err = tx.Stmt(s.queries[queryDeleteUser{}]).ExecContext(ctx, userId)

	if err != nil {
		return err
	}

	_, err = tx.Stmt(s.queries[queryReserveDeletedUserId{}]).ExecContext(ctx, userId)

	if err != nil {
		return err
	}

	return tx.Commit()
}
