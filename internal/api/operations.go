package api

import (
	"io"
	"net/http"
)

func (s *Service) registerNewUser(w http.ResponseWriter, r *http.Request) {

}

func (s *Service) getUserInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, "BEEP BOOP BEEP BEEP BOOP")
	w.WriteHeader(http.StatusOK)
}

func (s *Service) editUserInfo(w http.ResponseWriter, r *http.Request) {

}

func (s *Service) deleteUser(w http.ResponseWriter, r *http.Request) {

}
