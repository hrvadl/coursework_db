package server

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/go-chi/chi/v5"
	"github.com/hrvadl/coursework_db/pkg/controllers"
	"github.com/hrvadl/coursework_db/pkg/middleware"
	"github.com/hrvadl/coursework_db/pkg/repo"
	"go.uber.org/zap"
)

type Server struct {
	http.Server
	*HTTPServerArgs
}

type Controllers struct {
	Auth    *controllers.Auth
	Profile *controllers.Profile
	Deal    *controllers.Deal
}

type HTTPServerArgs struct {
	Session repo.Session
	Logger  *zap.SugaredLogger
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

	r.Use(middleware.WithCors())
	r.Use(middleware.WithHTTPLogger(s.Logger))

	r.With(middleware.RedirectAuthorized(s.Session)).Route("/auth", func(r chi.Router) {
		r.Get("/sign-up", s.Auth.ServeSignUpPage)
		r.Get("/sign-in", s.Auth.ServeSignInPage)
	})

	r.With(middleware.RedirectUnauthorized(s.Session)).Route("/", func(r chi.Router) {
		r.Get("/", s.Deal.ServeDealsPage)
		r.Get("/profile", s.Profile.ServeProfilePage)
	})

	// REST routes
	r.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.Post("/sign-up", s.Auth.HandleSignUp)
			r.Post("/sign-in", s.Auth.HandleSignIn)
		})

		r.Route("/profile", func(r chi.Router) {
		})

		r.Route("/stocks", func(r chi.Router) {
			r.Use(middleware.WithAuth(s.Session))
			r.Patch("/", func(w http.ResponseWriter, r *http.Request) {})
		})

		r.Route("/emitents", func(r chi.Router) {
			r.Use(middleware.WithAuth(s.Session))
			r.Patch("/", func(w http.ResponseWriter, r *http.Request) {
			})
		})

		r.Route("/securities", func(r chi.Router) {
			r.Use(middleware.WithAuth(s.Session))
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {})
		})

		r.Route("/deals", func(r chi.Router) {
			r.Get("/:ownerID", func(w http.ResponseWriter, r *http.Request) {})
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {})
			r.Patch("/:id", func(w http.ResponseWriter, r *http.Request) {})
			r.Delete("/:id", func(w http.ResponseWriter, r *http.Request) {})
		})

		r.Route("/transactions", func(r chi.Router) {
			r.Use(middleware.WithAuth(s.Session))
			r.Get("/:userID", func(w http.ResponseWriter, r *http.Request) {})
		})
	})

	fs := http.FileServer(http.Dir(filepath.Join("../templates")))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	return r
}
