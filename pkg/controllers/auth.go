package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/hrvadl/coursework_db/pkg/middleware"
	"github.com/hrvadl/coursework_db/pkg/models"
	"github.com/hrvadl/coursework_db/pkg/services"
	"github.com/hrvadl/coursework_db/pkg/templates"
)

func NewAuth(a services.Auth, t *templates.Resolver) *Auth {
	return &Auth{auth: a, t: t}
}

type Auth struct {
	t    *templates.Resolver
	auth services.Auth
}

func (a *Auth) ServeSignInPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	a.t.Execute(w, "sign-in.html", nil)
}

func (a *Auth) ServeSignUpPage(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	a.t.Execute(w, "sign-up.html", nil)
}

func (a *Auth) HandleSignIn(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	dto := &models.User{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}

	_, sess, err := a.auth.SignIn(dto)

	if err != nil {
		a.t.Execute(w, "sign-in-form", struct{ Error string }{Error: err.Error()})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     middleware.SessionCookie,
		Value:    strconv.FormatUint(uint64(sess.ID), 10),
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(time.Hour * 72),
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	})
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusOK)
}

func (a *Auth) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	dto := &models.User{
		Email:     r.FormValue("email"),
		Password:  r.FormValue("password"),
		FirstName: r.FormValue("first-name"),
		LastName:  r.FormValue("last-name"),
		Role:      r.FormValue("role"),
	}

	if _, err := a.auth.SignUp(dto); err != nil {
		a.t.Execute(w, "sign-up-form", templates.SignUpArgs{Error: err.Error()})
		return
	}

	a.t.Execute(w, "sign-up-form", templates.SignUpArgs{Success: true})
}

func (a *Auth) HandleLogOut(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     middleware.SessionCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	})

	http.Redirect(w, r, "/auth/sign-in", http.StatusMovedPermanently)
}
