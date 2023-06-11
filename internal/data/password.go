package data

import "context"

type Password string

const (
	queryPasswordVerifyName = "Password_Verify"
	queryPasswordVerify     = `
	SELECT true
	FROM users
	WHERE id = $1 AND password = $2;
	`
)

func (m *Password) Verify(userId string, s *Storage, ctx context.Context) error {
	return nil
}

// TODO
// func (m *Password) Set(userId string, s *Storage, ctx context.Context) error {
// 	return nil
// }
