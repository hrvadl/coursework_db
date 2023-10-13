package controllers

import (
	"net/http"

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
	http.ServeFile(w, r, a.t.ResolveHTML("sign-in.html"))
}

func (a *Auth) ServeSignUpPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, a.t.ResolveHTML("sign-up.html"))
}

func (u *Auth) HandleSignIn(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
}

func (u *Auth) HandleSignUp(w http.ResponseWriter, r *http.Request) {}
