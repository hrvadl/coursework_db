package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hrvadl/coursework_db/pkg/controllers"
	"github.com/hrvadl/coursework_db/pkg/middleware"
	"go.uber.org/zap"
)

type Server struct {
	http.Server
	*HTTPServerArgs
}

type Controllers struct {
	Auth *controllers.Auth
}

type HTTPServerArgs struct {
	Logger *zap.SugaredLogger
	*Controllers
}

func NewHTTP(h *HTTPServerArgs) *Server {
	return &Server{
		HTTPServerArgs: h,
	}
}

func (s *Server) ListenAndServe(port string) error {
	r := s.setupRoutes()
	return http.ListenAndServe(fmt.Sprintf(":%v", port), r)
}

func (s *Server) setupRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.NewCors())
	r.Use(middleware.NewHTTPLogger(s.Logger))

	r.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/sign-up", s.Auth.SignUp)
			r.Post("/sign-in", s.Auth.SignIn)
		})

		r.Route("/stocks", func(r chi.Router) {
			r.Patch("/", func(w http.ResponseWriter, r *http.Request) {})
		})

		r.Route("/emitents", func(r chi.Router) {
			r.Patch("/", func(w http.ResponseWriter, r *http.Request) {
			})
		})

		r.Route("/securities", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {})
		})

		r.Route("/deals", func(r chi.Router) {
			r.Get("/:ownerID", func(w http.ResponseWriter, r *http.Request) {})
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {})
			r.Patch("/:id", func(w http.ResponseWriter, r *http.Request) {})
			r.Delete("/:id", func(w http.ResponseWriter, r *http.Request) {})
		})

		r.Route("/transactions", func(r chi.Router) {
			r.Get("/:userID", func(w http.ResponseWriter, r *http.Request) {})
		})
	})

	return r
}
