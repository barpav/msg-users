package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/barpav/msg-users/internal/rest/models"
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

	var info *models.UserInfoV1
	info, err = s.storage.UserInfoV1(ctx, userId)

	if err != nil {
		logAndReturnErrorWithIssue(w, r, err, fmt.Sprintf("Failed to get user '%s' info before deletion.", userId))
		return
	}

	err = s.auth.EndAllSessions(ctx, userId)

	if err != nil {
		logAndReturnErrorWithIssue(w, r, err, fmt.Sprintf("Failed to end user '%s' sessions before deletion.", userId))
		return
	}

	err = s.storage.DeleteUser(ctx, userId)

	if err != nil {
		logAndReturnErrorWithIssue(w, r, err, fmt.Sprintf("Failed to delete user '%s'.", userId))
		return
	}

	log.Info().Msg(fmt.Sprintf("User '%s' deleted.", userId))

	if info.Picture != "" {
		go func() {
			err = s.fileStats.SendUsage(context.Background(), info.Picture, false)

			if err != nil {
				log.Err(err).Msg(fmt.Sprintf("Failed to send unused file '%s' statistics.", info.Picture))
			}
		}()
	}

	w.WriteHeader(http.StatusNoContent)
}
