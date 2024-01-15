package server

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	listenAddr string
	httpServer http.Server
}

func NewServer(listenAddr string) *Server {
	server := Server{listenAddr: listenAddr}
	router := mux.NewRouter()
	router.HandleFunc("/", server.Root)

	server.httpServer = http.Server{
		Addr:                         server.listenAddr,
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
	return &server
}

func (s *Server) StopServer() {
	ctx := context.Background();
	s.httpServer.Shutdown(ctx)
}

func (s *Server) Run() {
	log.Fatalln(s.httpServer.ListenAndServe())
}
