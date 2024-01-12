package server

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Server struct {
	listenAddr string
	httpServer http.Server
}

func NewServer(listenAddr string) *Server {
	return &Server{listenAddr: listenAddr}
}

func (s *Server) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/", s.Root)

	s.httpServer = http.Server{
		Addr:                         s.listenAddr,
		Handler:                      router,
		DisableGeneralOptionsHandler: false,
		TLSConfig:                    nil,
		ReadTimeout:                  0,
		ReadHeaderTimeout:            0,
		WriteTimeout:                 0,
		IdleTimeout:                  0,
		MaxHeaderBytes:               0,
		TLSNextProto:                 nil,
		ConnState:                    nil,
		ErrorLog:                     nil,
		BaseContext:                  nil,
		ConnContext:                  nil,
	}
	log.Fatalln(s.httpServer.ListenAndServe())
}
