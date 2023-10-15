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

func NewDeal(ds services.Deal, ss services.Security, t *templates.Resolver) *Deal {
	return &Deal{
		ds: ds,
		ss: ss,
		t:  t,
	}
}

type Deal struct {
	ds services.Deal
	ss services.Security
	t  *templates.Resolver
}

func (d *Deal) ServeDealsPage(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.Must(middleware.GetUserCtx(r.Context()))
	deals, err := d.ds.Get()

	if err != nil {
		d.t.Execute(w, "generic-error.html", templates.GenericErrorArgs{})
		return
	}

	securities, err := d.ss.Get()

	if err != nil {
		d.t.Execute(w, "generic-error.html", templates.GenericErrorArgs{})
		return
	}

	d.t.Execute(w, "deals.html", templates.DealsArgs{
		Deals:      deals,
		Role:       ctx.Role,
		Securities: securities,
	})
}

func (d *Deal) ServeDealPage(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	dealID, err := strconv.ParseInt(parts[len(parts)-1], 10, 64)

	if err != nil {
		d.t.Execute(w, "generic-error.html", templates.GenericErrorArgs{})
		return
	}

	deal, err := d.ds.GetByID(int(dealID))

	if err != nil || deal == nil {
		d.t.Execute(w, "generic-error.html", templates.GenericErrorArgs{})
		return
	}

	d.t.Execute(w, "deal.html", templates.DealArgs{Deal: deal})
}

func (d *Deal) HandleCreate(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.Must(middleware.GetUserCtx(r.Context()))
	r.ParseForm()

	securityID, err := strconv.ParseUint(r.FormValue("securityID"), 10, 64)

	if err != nil {
		d.t.Execute(w, "toast", templates.ToastArgs{Error: "Invalid security"})
		return
	}

	amount, err := strconv.ParseUint(r.FormValue("amount"), 10, 64)

	if err != nil {
		d.t.Execute(w, "toast", templates.ToastArgs{Error: "Invalid amount"})
		return
	}

	price, err := strconv.ParseFloat(r.FormValue("price"), 64)

	if err != nil {
		d.t.Execute(w, "toast", templates.ToastArgs{Error: "Invalid price"})
		return
	}

	dto := &models.Deal{
		OwnerID:    ctx.ID,
		SecurityID: uint(securityID),
		Amount:     uint(amount),
		Price:      price,
		Active:     true,
	}

	if _, err := d.ds.Create(dto); err != nil {
		d.t.Execute(w, "toast", templates.ToastArgs{Error: "Cannot create a deal"})
		return
	}

	w.Header().Set("HX-Trigger", "get-deals-event")
	d.t.Execute(w, "toast", templates.ToastArgs{Message: "Successfully created"})
}

func (d *Deal) HandleGet(w http.ResponseWriter, r *http.Request) {
	deals, _ := d.ds.Get()

	d.t.Execute(w, "deal-list", templates.DealListArgs{Deals: deals})
}
