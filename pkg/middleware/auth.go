package middleware

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/hrvadl/coursework_db/pkg/repo"
)

type UserCtx struct {
	ID   uint
	Role string
}

type key string

const User key = "user"

func WithAuth(session repo.Session) HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authH := r.Header.Get("Stock-Session-Auth")
			auth, err := strconv.ParseUint(authH, 10, 64)

			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			sess, err := session.GetByID(uint(auth))

			if err != nil || sess == nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if time.Now().Compare(sess.ValidUntil) != -1 {
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

func RedirectAuthorized(session repo.Session) HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authH := r.Header.Get("Stock-Session-Auth")
			auth, err := strconv.ParseUint(authH, 10, 64)

			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			sess, err := session.GetByID(uint(auth))

			if err != nil || sess == nil {
				next.ServeHTTP(w, r)
				return
			}

			if time.Now().Compare(sess.ValidUntil) != -1 {
				next.ServeHTTP(w, r)
				return
			}

			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		})
	}
}

func RedirectUnauthorized(session repo.Session) HTTPMiddleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authH := r.Header.Get("Stock-Session-Auth")
			auth, err := strconv.ParseUint(authH, 10, 64)

			if err != nil {
				http.Redirect(w, r, "/auth/sign-in", http.StatusMovedPermanently)
				return
			}

			sess, err := session.GetByID(uint(auth))

			if err != nil || sess == nil {
				http.Redirect(w, r, "/auth/sign-in", http.StatusMovedPermanently)
				return
			}

			if time.Now().Compare(sess.ValidUntil) != -1 {
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

func GetUserCtx(ctx context.Context) (*UserCtx, error) {
	val := ctx.Value(User)
	userCtx, ok := val.(UserCtx)

	if !ok {
		return nil, errors.New("cannot get user context")
	}

	return &userCtx, nil
}
