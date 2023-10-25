package controllers

import (
	"net/http"
	"strconv"

	"github.com/hrvadl/coursework_db/pkg/middleware"
	"github.com/hrvadl/coursework_db/pkg/models"
	"github.com/hrvadl/coursework_db/pkg/services"
	"github.com/hrvadl/coursework_db/pkg/templates"
)

func NewProfile(
	us services.User,
	ds services.Deal,
	ts services.Transaction,
	ses services.Security,
	t *templates.Resolver) *Profile {
	return &Profile{
		us:  us,
		ds:  ds,
		ts:  ts,
		ses: ses,
		t:   t,
	}
}

type Profile struct {
	us  services.User
	ds  services.Deal
	ts  services.Transaction
	ses services.Security
	t   *templates.Resolver
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

	profile, err := p.us.GetByID(int(userCtx.ID))

	if err != nil {
		p.t.Execute(w, "generic-error.html", templates.GenericErrorArgs{})
		return
	}

	transactions, err := p.ts.Get(int(profile.ID))

	if err != nil {
		p.t.Execute(w, "generic-error.html", templates.GenericErrorArgs{})
		return
	}

	securities, err := p.ses.Get()

	if err != nil {
		p.t.Execute(w, "generic-error.html", templates.GenericErrorArgs{})
		return
	}

	p.t.Execute(w, "profile.html", templates.ProfileArgs{
		User:         profile,
		Logined:      true,
		Transactions: transactions,
		Securities:   securities,
	})
}

func (p *Profile) HandlePatch(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.Must(
		middleware.GetUserCtx(r.Context()),
	)

	r.ParseForm()
	amount, err := strconv.ParseInt(r.FormValue("amount"), 10, 64)

	if err != nil {
		p.t.Execute(w, "toast", templates.ToastArgs{Error: "Invalid money amount"})
		return
	}

	isWithdraw, _ := strconv.ParseBool(r.FormValue("withdraw"))

	switch isWithdraw {
	case true:
		err = p.us.WithdrawMoney(int(userCtx.ID), int(amount))
	default:
		err = p.us.AddMoney(int(userCtx.ID), int(amount))
	}

	if err != nil {
		p.t.Execute(w, "toast", templates.ToastArgs{Error: "Cannot update the user"})
		return
	}

	w.Header().Set("HX-Trigger", "refresh-general-info")
	p.t.Execute(w, "toast", templates.ToastArgs{Message: "Successfully updated the user"})
}

func (p *Profile) HandleGetGeneralInfo(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.Must(
		middleware.GetUserCtx(r.Context()),
	)

	profile, err := p.us.GetByID(int(userCtx.ID))

	if err != nil {
		p.t.Execute(w, "toast", templates.ToastArgs{Error: "User not found"})
		return
	}

	p.t.Execute(w, "general-info", templates.GeneralProfileInfoArgs{
		User: profile,
	})
}
