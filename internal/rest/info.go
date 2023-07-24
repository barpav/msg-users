package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const mimeTypeUserInfoV1 = "application/vnd.userInfo.v1+json"

// https://barpav.github.io/msg-api-spec/#/users/get_users
func (s *Service) getUserInfo(w http.ResponseWriter, r *http.Request) {
	switch r.Header.Get("Accept") {
	case "", mimeTypeUserInfoV1: // including if not specified
		s.getUserInfoV1(w, r)
	default:
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
}

func (s *Service) getUserInfoV1(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("id")

	if userId == "" {
		userId = authenticatedUser(r)
	}

	info, err := s.storage.UserInfoV1(r.Context(), userId)

	if err == nil {
		if info == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", mimeTypeUserInfoV1)
		err = json.NewEncoder(w).Encode(info)
	}

	if err != nil {
		logAndReturnErrorWithIssue(w, r, err, fmt.Sprintf("Failed to get user '%s' info.", userId))
		return
	}

	w.WriteHeader(http.StatusOK)
}
