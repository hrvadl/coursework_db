package controllers

import (
	"net/http"

	"github.com/hrvadl/coursework_db/pkg/services"
)

func NewUser(a services.Auth) *Auth {
	return &Auth{auth: a}
}

type Auth struct {
	auth services.Auth
}

func (u *Auth) SignIn(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/html")
}

func (u *Auth) SignUp(w http.ResponseWriter, r *http.Request) {}
