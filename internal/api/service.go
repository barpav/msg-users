package api

import (
	"net/http"

	"github.com/barpav/msg-users/internal/data"
)

type Service struct {
	storage *data.Storage
}

func (s *Service) Init() error {
	s.storage = &data.Storage{}
	return s.storage.Init()
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// URL handling is reverse proxy's concern
	switch r.Method {
	case http.MethodPost:
		s.registerNewUser(w, r)
	case http.MethodGet:
		s.getUserInfo(w, r)
	case http.MethodPut:
		s.editUserInfo(w, r)
	case http.MethodDelete:
		s.deleteUser(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
