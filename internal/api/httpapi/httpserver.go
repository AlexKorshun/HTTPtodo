package httpapi

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	port    string
	handler http.Handler
}

func NewServer(port string, handler http.Handler) Server {
	return Server{port: port, handler: handler}
}

func (s *Server) Run() error {
	log.Printf("сервер запущен на :%s\n", s.port)

	return http.ListenAndServe(fmt.Sprintf(":%s", s.port), s.handler)
}
