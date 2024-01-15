package server

import (
	"net/http"
)

func (s *Server) Root(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello there\n"))
}
