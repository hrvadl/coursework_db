package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/hrvadl/coursework_db/pkg/models"
	"github.com/hrvadl/coursework_db/pkg/repo"
	"github.com/hrvadl/coursework_db/pkg/templates"
)

type UserCtx struct {
	ID   uint
	Role string
}

type key string

const (
	User          key    = "user"
	SessionCookie string = "Stock-Session-Auth"
)

func NewAuth(s repo.Session, t *templates.Resolver) *Auth {
	return &Auth{
		session: s,
		templ:   t,
	}
}

type Auth struct {
	session repo.Session
	templ   *templates.Resolver
}

// Only extracts the session cookie (if it exists) and puts it in the context
// Not panic and do nothing in case of session cookie is not present
func (m *Auth) WithUserCredsExtractor() HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess, err := m.getAuthCredsFromRequest(r)

			if err == nil && sess != nil {
				ctx := context.WithValue(r.Context(), User, UserCtx{
					ID:   sess.UserID,
					Role: sess.UserRole,
				})

				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Extracts the session cookie and puts it in the context
// In case of session cookie is not present return 401
func (m *Auth) WithAuth() HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess, err := m.getAuthCredsFromRequest(r)

			if err != nil || sess == nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), User, UserCtx{
				ID:   sess.UserID,
				Role: sess.UserRole,
			})

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func (m *Auth) WithoutAuth() HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess, err := m.getAuthCredsFromRequest(r)

			if err == nil && sess != nil {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Extracts the session cookie and redirects to the home page if it is present
func (m *Auth) RedirectAuthorized() HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess, err := m.getAuthCredsFromRequest(r)

			if err != nil || sess == nil {
				next.ServeHTTP(w, r)
				return
			}

			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		})
	}
}

// Extracts the session cookie and puts it in the context
// In case of session cookie is not present redirects to the auth page
func (m *Auth) RedirectUnauthorized() HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess, err := m.getAuthCredsFromRequest(r)

			if err != nil {
				http.Redirect(w, r, "/auth/sign-in", http.StatusMovedPermanently)
				return
			}

			ctx := context.WithValue(r.Context(), User, UserCtx{
				ID:   sess.UserID,
				Role: sess.UserRole,
			})

			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func (m *Auth) WithSameUserID() HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			parts := strings.Split(r.URL.Path, "/")

			userID, err := strconv.ParseInt(parts[len(parts)-1], 10, 64)
			if err != nil {
				m.templ.Execute(w, "generic-error.html", templates.GenericErrorArgs{})
				return
			}

			userCtx, err := GetUserCtx(r.Context())

			if err != nil {
				m.templ.Execute(w, "toast", templates.ToastArgs{Error: "Something went wrong"})
				return
			}

			if userCtx.ID != uint(userID) {
				m.templ.Execute(w, "toast", templates.ToastArgs{Error: "You can update only yours profile"})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (m *Auth) getAuthCredsFromRequest(r *http.Request) (*models.Session, error) {
	authCookie, err := r.Cookie(SessionCookie)

	if err != nil {
		return nil, err
	}

	auth, err := strconv.ParseUint(authCookie.Value, 10, 64)

	if err != nil {
		return nil, err
	}

	sess, err := m.session.GetByID(uint(auth))
	return sess, err
}
