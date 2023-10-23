package server

import (
	"fmt"
	"net/http"

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
	Auth      *controllers.Auth
	Profile   *controllers.Profile
	Deal      *controllers.Deal
	Inventory *controllers.Inventory
}

type HTTPServerArgs struct {
	AuthM   *middleware.Auth
	CorsM   *middleware.Cors
	LoggerM *middleware.HTTPLogger

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

	r.Use(s.CorsM.WithCors())
	r.Use(s.LoggerM.WithHTTPLogger())

	r.With(s.AuthM.RedirectAuthorized()).Route("/auth", func(r chi.Router) {
		r.Get("/sign-up", s.Auth.ServeSignUpPage)
		r.Get("/sign-in", s.Auth.ServeSignInPage)
	})

	r.With(s.AuthM.WithUserCredsExtractor()).Group(func(r chi.Router) {
		r.Get("/", s.Deal.ServeDealsPage)
		r.Get("/deals/{id}", s.Deal.ServeDealPage)
	})

	r.With(s.AuthM.RedirectUnauthorized()).Group(func(r chi.Router) {
		r.Get("/profile", s.Profile.ServeProfilePage)
	})

	r.Route("/v1", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.With(s.AuthM.WithoutAuth()).Group(func(r chi.Router) {
				r.Post("/sign-up", s.Auth.HandleSignUp)
				r.Post("/sign-in", s.Auth.HandleSignIn)
			})

			r.With(s.AuthM.WithAuth()).Get("/log-out", s.Auth.HandleLogOut)
		})

		r.With(s.AuthM.WithAuth()).Group(func(r chi.Router) {
			r.Route("/profile", func(r chi.Router) {
				r.With(s.AuthM.WithSameUserID()).Group(func(r chi.Router) {
					r.Patch("/{id}", s.Profile.HandlePatch)
					r.Get("/general-info/{id}", s.Profile.HandleGetGeneralInfo)
				})
			})

			r.With(s.AuthM.WithSameUserID()).Route("/inventory", func(r chi.Router) {
				r.Patch("/{userID}", s.Inventory.HandlePatch)
				r.Get("/{userID}", s.Inventory.HandleGetInventoryInfo)
			})

			r.Route("/deals", func(r chi.Router) {
				r.Get("/", s.Deal.HandleGet)
				r.Post("/", s.Deal.HandleCreate)
				r.Patch("/{id}", func(w http.ResponseWriter, r *http.Request) {})
				r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {})
			})

			r.Route("/transactions", func(r chi.Router) {
				r.Get("/{userID}", func(w http.ResponseWriter, r *http.Request) {})
			})
		})

	})
	return r
}
