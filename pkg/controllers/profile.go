package controllers

import (
	"net/http"

	"github.com/hrvadl/coursework_db/pkg/middleware"
	"github.com/hrvadl/coursework_db/pkg/models"
	"github.com/hrvadl/coursework_db/pkg/services"
	"github.com/hrvadl/coursework_db/pkg/templates"
)

func NewProfile(
	es services.Emitent,
	ss services.Stock,
	ds services.Deal,
	t *templates.Resolver) *Profile {
	return &Profile{
		es: es,
		ss: ss,
		ds: ds,
		t:  t,
	}
}

type Profile struct {
	es services.Emitent
	ss services.Stock
	ds services.Deal
	t  *templates.Resolver
}

type ProfileStrategy interface {
	GetByID(id int) (*models.User, error)
}

func (p *Profile) ServeProfilePage(w http.ResponseWriter, r *http.Request) {
	userCtx, err := middleware.GetUserCtx(r.Context())

	if err != nil {
		// TODO: error page
		return
	}

	var profileStrategy ProfileStrategy
	switch userCtx.Role {
	case models.EmitentRole:
		profileStrategy = p.es
	case models.StockRole:
		profileStrategy = p.ss
	}

	_, err = profileStrategy.GetByID(int(userCtx.ID))

	if err != nil {
		// TODO: error page
		return
	}

	http.ServeFile(w, r, p.t.ResolveHTML("profile.html"))
}
