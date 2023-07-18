package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

func (s *Service) traceInternalServerError(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				err := errors.New(fmt.Sprintf("Recovered from panic: %v", rec))
				log.Err(err).Msg(fmt.Sprintf("Internal server error (issue: %s).", requestId(r)))

				addIssueHeader(w, r)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func addIssueHeader(w http.ResponseWriter, r *http.Request) {
	w.Header()["issue"] = []string{requestId(r)} // lowercase - non-canonical (vendor) header
}

func requestId(r *http.Request) string {
	id := r.Header.Get("request-id") // set by api-gateway

	if id != "" {
		return id
	}

	return "untraced"
}
