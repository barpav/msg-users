package api

import (
	"io"
	"net/http"
)

func (s *Service) getUserInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, "BEEP BOOP BEEP BEEP BOOP")
	w.WriteHeader(http.StatusOK)
}
