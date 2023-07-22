package models

import (
	"encoding/json"
	"errors"
	"io"
)

// Schema: userProfilePassword.v1
type UserProfilePasswordV1 struct {
	Current string
	New     string
}

func (m *UserProfilePasswordV1) Deserialize(data io.Reader) error {
	if json.NewDecoder(data).Decode(m) != nil {
		return errors.New("Profile data violates 'userProfilePassword.v1' schema.")
	}

	return validatePassword(m.New)
}
