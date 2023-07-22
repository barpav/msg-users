package models

import (
	"errors"
	"unicode"
)

func validatePassword(password string) (err error) {
	type checklist struct {
		length    bool
		uppercase bool
		lowercase bool
		digit     bool
	}

	check := checklist{length: len([]rune(password)) >= 8}

	for _, ch := range password {
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
