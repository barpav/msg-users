package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

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

type authenticatedUserId struct{}

func (s *Service) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			next.ServeHTTP(w, r) // authentication not required (user registration)
		case http.MethodDelete:
			s.authenticateByCredentials(w, r, next)
		default:
			s.authenticateBySessionKey(w, r, next)
		}
	})
}

func (s *Service) authenticateBySessionKey(w http.ResponseWriter, r *http.Request, next http.Handler) {
	var userId string
	var err error
	sessionKey := sessionKey(r)

	if sessionKey != "" {
		userId, err = s.auth.ValidateSession(r.Context(), sessionKey, userIP(r), userAgent(r))
	}

	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("Session key authentication failed (issue: %s).", requestId(r)))

		addIssueHeader(w, r)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if userId == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx := context.WithValue(r.Context(), authenticatedUserId{}, userId)

	next.ServeHTTP(w, r.WithContext(ctx))
}

func (s *Service) authenticateByCredentials(w http.ResponseWriter, r *http.Request, next http.Handler) {
	var authenticated bool
	var err error

	userId, password, parsed := r.BasicAuth()

	if parsed {
		authenticated, err = s.storage.ValidateCredentials(r.Context(), userId, password)
	}

	if err != nil {
		log.Err(err).Msg(fmt.Sprintf("Credentials authentication failed (issue: %s).", requestId(r)))

		addIssueHeader(w, r)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !authenticated {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	ctx := context.WithValue(r.Context(), authenticatedUserId{}, userId)

	next.ServeHTTP(w, r.WithContext(ctx))
}

func sessionKey(r *http.Request) string {
	parts := strings.Split(r.Header.Get("Authorization"), "Bearer")

	if len(parts) != 2 {
		return ""
	}

	return strings.TrimSpace(parts[1])
}

func authenticatedUser(r *http.Request) (id string) {
	id, _ = r.Context().Value(authenticatedUserId{}).(string)
	return id
}

func userIP(r *http.Request) string {
	ip := r.Header.Get("remote-addr") // set by api-gateway

	if ip != "" {
		return ip
	}

	return "unknown"
}

func userAgent(r *http.Request) string {
	agent := r.Header.Get("User-Agent")

	if agent != "" {
		return agent
	}

	return "unknown"
}
