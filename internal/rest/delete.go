package rest

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

// https://barpav.github.io/msg-api-spec/#/users/delete_users
func (s *Service) deleteUser(w http.ResponseWriter, r *http.Request) {
	confirmCode := r.URL.Query().Get("confirmation-code")

	userId, ctx := authenticatedUser(r), r.Context()
	var err error

	if confirmCode == "" {
		confirmCode, err = s.storage.GenerateUserDeletionCode(ctx, userId)

		if err != nil {
			logAndReturnErrorWithIssue(w, r, err,
				fmt.Sprintf("Failed to generate user '%s' deletion confirmation code.", userId))
			return
		}

		w.Header()["confirmation-code"] = []string{confirmCode} // lowercase - non-canonical (vendor) header
		w.WriteHeader(http.StatusAccepted)
		return
	}
	var valid bool
	valid, err = s.storage.ValidateUserDeletionCode(ctx, userId, confirmCode)

	if err != nil {
		logAndReturnErrorWithIssue(w, r, err,
			fmt.Sprintf("Failed to validate user '%s' deletion confirmation code.", userId))
		return
	}

	if !valid {
		w.WriteHeader(http.StatusExpectationFailed)
		return
	}

	err = s.storage.DeleteUser(ctx, userId)

	if err != nil {
		logAndReturnErrorWithIssue(w, r, err, fmt.Sprintf("Failed to delete user '%s'.", userId))
		return
	}

	log.Info().Msg(fmt.Sprintf("User '%s' deleted.", userId))

	w.WriteHeader(http.StatusNoContent)
}
