package models

import (
	"encoding/json"
	"errors"
	"io"
)

// Schema: userProfileCommon.v1
type UserProfileCommonV1 struct {
	Name    string
	Picture string
}

func (m *UserProfileCommonV1) Deserialize(data io.Reader) error {
	if json.NewDecoder(data).Decode(m) != nil {
		return errors.New("Profile data violates 'userProfileCommon.v1' schema.")
	}

	return m.validate()
}

func (m *UserProfileCommonV1) validate() (err error) {
	if len([]rune(m.Name)) > 150 {
		err = errors.Join(err, errors.New("User name must be between 1 and 150 characters."))
	}

	if m.Picture != "" && len(m.Picture) != 24 {
		err = errors.Join(err, errors.New("User profile picture must be 24-character file id."))
	}

	return err
}
