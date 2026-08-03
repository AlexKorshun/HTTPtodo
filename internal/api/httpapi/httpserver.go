package httpapi

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	port string
	mux  *http.ServeMux
}

func NewServer(port string, mux *http.ServeMux) Server {
	return Server{port: port, mux: mux}
}

func (s *Server) Run() error {
	log.Printf("сервер запущен на :%s\n", s.port)

	return http.ListenAndServe(fmt.Sprintf(":%s", s.port), s.mux)
}
