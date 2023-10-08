package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	http.Server
}

func NewHTTP() *Server {
	return &Server{}
}

func (s *Server) ListenAndServe(port string) error {
	r := chi.NewRouter()

	r.Use(NewCorsMiddleware())

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	})

	return http.ListenAndServe(port, r)
}
