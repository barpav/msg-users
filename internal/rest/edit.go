package rest

import (
	"fmt"
	"net/http"

	"github.com/barpav/msg-users/internal/rest/models"
	"github.com/rs/zerolog/log"
)

// https://barpav.github.io/msg-api-spec/#/users/patch_users
func (s *Service) editUserInfo(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("Content-Type") {
	case "application/vnd.userProfileCommon.v1+json":
		s.editCommonProfileInfoV1(w, r)
	case "application/vnd.userProfilePassword.v1+json":
		s.changePasswordV1(w, r)
	default:
		http.Error(w, "Unsupported profile data (invalid media type).", 400)
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
		log.Err(err).Msg(fmt.Sprintf("Failed to update user '%s' common profile info (issue: %s).",
			userId, requestId(r)))

		addIssueHeader(w, r)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) changePasswordV1(w http.ResponseWriter, r *http.Request) {

}
