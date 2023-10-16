package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/hrvadl/coursework_db/pkg/middleware"
	"github.com/hrvadl/coursework_db/pkg/models"
	"github.com/hrvadl/coursework_db/pkg/services"
	"github.com/hrvadl/coursework_db/pkg/templates"
)

func NewProfile(
	es services.Emitent,
	ss services.Stock,
	ds services.Deal,
	ts services.Transaction,
	t *templates.Resolver) *Profile {
	return &Profile{
		es: es,
		ss: ss,
		ds: ds,
		ts: ts,
		t:  t,
	}
}

type Profile struct {
	es services.Emitent
	ss services.Stock
	ds services.Deal
	ts services.Transaction
	t  *templates.Resolver
}

type ProfileStrategy interface {
	GetByID(id int) (*models.User, error)
	Patch(user *models.User) (*models.User, error)
}

func (p *Profile) ServeProfilePage(w http.ResponseWriter, r *http.Request) {
	userCtx, err := middleware.GetUserCtx(r.Context())

	if err != nil {
		p.t.Execute(w, "generic-error.html", templates.GenericErrorArgs{})
		return
	}

	profileStrategy, err := p.chooseUserStrategy(userCtx)

	if err != nil {
		p.t.Execute(w, "generic-error.html", templates.GenericErrorArgs{})
		return
	}

	profile, err := profileStrategy.GetByID(int(userCtx.ID))

	if err != nil {
		p.t.Execute(w, "generic-error.html", templates.GenericErrorArgs{})
		return
	}

	transactions, err := p.ts.Get(int(profile.ID))

	if err != nil {
		p.t.Execute(w, "generic-error.html", templates.GenericErrorArgs{})
		return
	}

	p.t.Execute(w, "profile.html", templates.ProfileArgs{
		User:         profile,
		Logined:      true,
		Transactions: transactions,
	})
}

func (p *Profile) HandlePatch(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	userID, err := strconv.ParseInt(parts[len(parts)-1], 10, 64)
	if err != nil {
		p.t.Execute(w, "generic-error.html", templates.GenericErrorArgs{})
		return
	}

	userCtx, err := middleware.GetUserCtx(r.Context())

	if err != nil {
		p.t.Execute(w, "toast", templates.ToastArgs{Error: "Something went wrong"})
		return
	}

	if userCtx.ID != uint(userID) {
		p.t.Execute(w, "toast", templates.ToastArgs{Error: "You can update only yours profile"})
		return
	}

	profileStrategy, err := p.chooseUserStrategy(userCtx)

	if err != nil {
		p.t.Execute(w, "toast", templates.ToastArgs{Error: "Something went wrong"})
		return
	}

	profile, err := profileStrategy.GetByID(int(userID))

	if err != nil {
		p.t.Execute(w, "toast", templates.ToastArgs{Error: "User not found"})
		return
	}

	r.ParseForm()

	amount, err := strconv.ParseInt(r.FormValue("amount"), 10, 64)

	if err != nil {
		p.t.Execute(w, "toast", templates.ToastArgs{Error: "Invalid money amount"})
		return
	}

	dto := &models.User{
		ID:      uint(userID),
		Balance: profile.Balance - int(amount),
	}

	if _, err := profileStrategy.Patch(dto); err != nil {
		p.t.Execute(w, "toast", templates.ToastArgs{Error: "Cannot update the user"})
		return
	}

	w.Header().Set("HX-Trigger", "refresh-general-info")
	p.t.Execute(w, "toast", templates.ToastArgs{Message: "Successfully updated the user"})
}

func (p *Profile) chooseUserStrategy(userCtx *middleware.UserCtx) (ProfileStrategy, error) {
	var profileStrategy ProfileStrategy
	switch userCtx.Role {
	case models.EmitentRole:
		profileStrategy = p.es
	case models.StockRole:
		profileStrategy = p.ss
	}

	return profileStrategy, nil
}
