package middleware

import (
	"github.com/go-chi/cors"
)

func NewCors() *Cors {
	return &Cors{}
}

type Cors struct{}

func (m *Cors) WithCors() HTTPMiddleware {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	})
}
