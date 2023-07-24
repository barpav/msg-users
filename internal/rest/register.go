package rest

import (
	"fmt"
	"net/http"

	"github.com/barpav/msg-users/internal/rest/models"
	"github.com/rs/zerolog/log"
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

type ErrUserAlreadyExists interface {
	Error() string
	ImplementsUserAlreadyExistsError()
}

func (s *Service) registerNewUserV1(w http.ResponseWriter, r *http.Request) {
	userInfo := models.NewUserV1{}
	err := userInfo.Deserialize(r.Body)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	err = s.storage.CreateUser(r.Context(), userInfo.Id, userInfo.Name, userInfo.Password)

	if err != nil {
		if _, ok := err.(ErrUserAlreadyExists); ok {
			w.WriteHeader(http.StatusConflict)
			return
		}

		logAndReturnErrorWithIssue(w, r, err, "User registration failed")
		return
	}

	log.Info().Msg(fmt.Sprintf("User '%s' successfully registered.", userInfo.Id))

	w.WriteHeader(http.StatusCreated)
}
