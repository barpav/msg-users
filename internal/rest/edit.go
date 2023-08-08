package rest

import (
	"fmt"
	"net/http"

	"github.com/barpav/msg-users/internal/rest/models"
)

// https://barpav.github.io/msg-api-spec/#/users/patch_users
func (s *Service) editUserInfo(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("Content-Type") {
	case "application/vnd.userProfileCommon.v1+json":
		s.editCommonProfileInfoV1(w, r)
	case "application/vnd.userProfilePassword.v1+json":
		s.changePasswordV1(w, r)
	default:
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
}

func (s *Service) editCommonProfileInfoV1(w http.ResponseWriter, r *http.Request) {
	info := models.UserProfileCommonV1{}
	err := info.Deserialize(r.Body)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	userId := authenticatedUser(r)

	err = s.storage.UpdateCommonProfileInfoV1(r.Context(), userId, &info)

	if err != nil {
		logAndReturnErrorWithIssue(w, r, err,
			fmt.Sprintf("Failed to update user '%s' common profile info.", userId))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) changePasswordV1(w http.ResponseWriter, r *http.Request) {
	passwords := models.UserProfilePasswordV1{}
	err := passwords.Deserialize(r.Body)

	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	userId, ctx := authenticatedUser(r), r.Context()
	var currentIsValid bool
	currentIsValid, err = s.storage.ValidateCredentials(ctx, userId, passwords.Current)

	if err == nil {
		if !currentIsValid {
			http.Error(w, "Invalid current user password.", 400)
			return
		}

		err = s.storage.ChangePassword(ctx, userId, passwords.New)
	}

	if err != nil {
		logAndReturnErrorWithIssue(w, r, err, fmt.Sprintf("Failed to change user '%s' password.", userId))
		return
	}

	w.WriteHeader(http.StatusOK)
}
