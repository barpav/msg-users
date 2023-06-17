package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode"

	"github.com/rs/zerolog/log"

	"github.com/barpav/msg-users/internal/data"
)

// https://barpav.github.io/msg-api-spec/#/users/post_users
func (s *Service) registerNewUser(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("Content-Type") {
	case "application/vnd.newUser.v1+json":
		s.registerNewUserV1(w, r)
	default:
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
}

func (s *Service) registerNewUserV1(w http.ResponseWriter, r *http.Request) {
	userInfo := newUserV1{}
	err := userInfo.deserialize(r.Body)

	if err != nil {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		http.Error(w, err.Error(), 400)
		return
	}

	user := data.NewUser{
		Id:       userInfo.Id,
		Name:     userInfo.Name,
		Password: userInfo.Password,
	}

	err, exists := user.Create(s.storage, r.Context())

	if exists {
		w.WriteHeader(http.StatusConflict)
		return
	}

	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("User registration failed (issue: %s).", r.Header.Get("request-id")))

		w.Header()["issue"] = []string{r.Header.Get("request-id")} // lowercase - non-canonical (vendor) header
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Info().Msg(fmt.Sprintf("User '%s' successfully registered.", user.Id))

	w.WriteHeader(http.StatusCreated)
}

// schema: newUser.v1
type newUserV1 struct {
	Id       string
	Name     string
	Password string
}

func (m *newUserV1) deserialize(data io.ReadCloser) error {
	if json.NewDecoder(data).Decode(m) != nil {
		return errors.New("New user data violates 'newUser.v1' schema.")
	}

	m.Id = strings.ToLower(m.Id)
	m.Id = strings.TrimSpace(m.Id)

	m.Name = strings.TrimSpace(m.Name)

	return m.validate()
}

func (m *newUserV1) validate() (err error) {
	err = errors.Join(err, m.validateId())

	if nLen := len([]rune(m.Name)); nLen < 1 || nLen > 150 {
		err = errors.Join(err, errors.New("User name must be between 1 and 150 characters."))
	}

	err = errors.Join(err, m.validatePassword())

	return err
}

func (m *newUserV1) validateId() (err error) {
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

func (m *newUserV1) validatePassword() (err error) {
	type checklist struct {
		length    bool
		uppercase bool
		lowercase bool
		digit     bool
	}

	check := checklist{length: len([]rune(m.Password)) >= 8}

	for _, ch := range m.Password {
		switch {
		case unicode.IsUpper(ch):
			check.uppercase = true
		case unicode.IsLower(ch):
			check.lowercase = true
		case unicode.IsDigit(ch):
			check.digit = true
		}
	}

	if !check.length {
		err = errors.Join(err, errors.New("User password must be at least 8 characters long."))
	}

	if !check.uppercase {
		err = errors.Join(err, errors.New("User password must contain at least one uppercase letter."))
	}

	if !check.lowercase {
		err = errors.Join(err, errors.New("User password must contain at least one lowercase letter."))
	}

	if !check.digit {
		err = errors.Join(err, errors.New("User password must contain at least one digit."))
	}

	return err
}
