package models

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
)

// Schema: newUser.v1
type NewUserV1 struct {
	Id       string
	Name     string
	Password string
}

func (m *NewUserV1) Deserialize(data io.ReadCloser) error {
	if json.NewDecoder(data).Decode(m) != nil {
		return errors.New("New user data violates 'newUser.v1' schema.")
	}

	m.Id = strings.ToLower(m.Id)
	m.Id = strings.TrimSpace(m.Id)

	m.Name = strings.TrimSpace(m.Name)

	return m.validate()
}

func (m *NewUserV1) validate() (err error) {
	err = errors.Join(err, m.validateId())

	if nLen := len([]rune(m.Name)); nLen < 1 || nLen > 150 {
		err = errors.Join(err, errors.New("User name must be between 1 and 150 characters."))
	}

	err = errors.Join(err, validatePassword(m.Password))

	return err
}

func (m *NewUserV1) validateId() (err error) {
	for _, ch := range m.Id {
		if (ch < 'a' || ch > 'z') && (ch < '0' || ch > '9') && ch != '_' && ch != '-' {
			err = errors.Join(err, errors.New("User id can contain only 'a-z', '0-9', '-' and '_' characters."))
			break
		}
	}

	if l := len([]rune(m.Id)); l < 1 || l > 50 {
		err = errors.Join(err, errors.New("User id must be between 1 and 50 characters."))
	}

	return err
}
