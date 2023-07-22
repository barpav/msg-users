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
	if nLen := len([]rune(m.Name)); nLen < 1 || nLen > 150 {
		err = errors.New("User name must be between 1 and 150 characters.")
	}

	return err
}
