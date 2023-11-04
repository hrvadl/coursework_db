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

func NewDeal(ds services.Deal, ss services.Security, is services.Inventory, t *templates.Resolver) *Deal {
	return &Deal{
		ds: ds,
		ss: ss,
		is: is,
		t:  t,
	}
}

type Deal struct {
	ds services.Deal
	ss services.Security
	is services.Inventory
	t  *templates.Resolver
}

func (d *Deal) ServeDealsPage(w http.ResponseWriter, r *http.Request) {
	deals, err := d.ds.Get()

	if err != nil {
		d.t.Execute(w, "generic-error.html", templates.GenericErrorArgs{Logined: true})
		return
	}

	securities, err := d.ss.Get()

	if err != nil {
		d.t.Execute(w, "generic-error.html", templates.GenericErrorArgs{Logined: true})
		return
	}

	var role string
	ctx, _ := middleware.GetUserCtx(r.Context())
	if ctx == nil {
		role = ""
	} else {
		role = ctx.Role
	}

	w.WriteHeader(http.StatusOK)
	d.t.Execute(w, "deals.html", templates.DealsArgs{
		Deals:      deals,
		Role:       role,
		Securities: securities,
		Logined:    ctx != nil,
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

	ctx, _ := middleware.GetUserCtx(r.Context())
	if ctx == nil {
		d.t.Execute(w, "deal.html", templates.DealArgs{
			Deal:    deal,
			Logined: false,
		})
		return
	}

	has, _ := d.is.GetUserInventoryBySecurityID(int(ctx.ID), int(deal.SecurityID))
	if err != nil {
		d.t.Execute(w, "generic-error.html", templates.GenericErrorArgs{})
		return
	}

	var amountHas int
	switch has {
	case nil:
		amountHas = 0
	default:
		amountHas = int(has.Amount)
	}

	w.WriteHeader(http.StatusOK)
	d.t.Execute(w, "deal.html", templates.DealArgs{
		Deal:      deal,
		Logined:   true,
		AmountHas: amountHas,
		IsOwner:   ctx.ID == deal.OwnerID,
	})
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
		Sell:       r.FormValue("type") == "sell",
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

func (d *Deal) HandlePatch(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.Must(
		middleware.GetUserCtx(r.Context()),
	)

	parts := strings.Split(r.URL.Path, "/")
	dealID, err := strconv.ParseInt(parts[len(parts)-1], 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()

	amount, err := strconv.ParseInt(r.FormValue("amount"), 10, 64)

	if err != nil {
		d.t.Execute(w, "toast", templates.ToastArgs{Error: "Invalid amount"})
		return
	}

	price, err := strconv.ParseFloat(r.FormValue("price"), 64)

	if err != nil {
		d.t.Execute(w, "toast", templates.ToastArgs{Error: "Invalid price"})
		return
	}

	deal, err := d.ds.GetByID(int(dealID))

	if err != nil {
		d.t.Execute(w, "toast", templates.ToastArgs{Error: "Deal does not exist"})
		return
	}

	if deal.OwnerID != userCtx.ID {
		d.t.Execute(w, "toast", templates.ToastArgs{Error: "You can update only deals you own"})
		return
	}

	_, err = d.ds.Patch(&models.Deal{
		ID:     uint(dealID),
		Amount: uint(amount),
		Price:  price,
	})

	if err != nil {
		d.t.Execute(w, "toast", templates.ToastArgs{Error: "Failed to update the deal"})
		return
	}

	w.Header().Set("HX-Trigger", "get-general-info")
	d.t.Execute(w, "toast", templates.ToastArgs{Message: "Successfully updated the deal"})
}

func (d *Deal) HandleDelete(w http.ResponseWriter, r *http.Request) {
	userCtx := middleware.Must(
		middleware.GetUserCtx(r.Context()),
	)

	parts := strings.Split(r.URL.Path, "/")
	dealID, err := strconv.ParseInt(parts[len(parts)-1], 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	deal, err := d.ds.GetByID(int(dealID))

	if err != nil {
		d.t.Execute(w, "toast", templates.ToastArgs{Error: "Deal does not exist"})
		return
	}

	if deal.OwnerID != userCtx.ID {
		d.t.Execute(w, "toast", templates.ToastArgs{Error: "You can delete only deals you own"})
		return
	}

	if err = d.ds.Delete(int(dealID)); err != nil {
		d.t.Execute(w, "toast", templates.ToastArgs{Error: "Failed to delete the deal"})
		return
	}

	w.Header().Set("HX-Redirect", "/")
	d.t.Execute(w, "toast", templates.ToastArgs{Message: "Successfully deleted the deal"})
}

func (d *Deal) HandleGetGeneralInfo(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	dealID, err := strconv.ParseInt(parts[len(parts)-1], 10, 64)

	if err != nil {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	deal, err := d.ds.GetAllByID(int(dealID))

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	d.t.Execute(w, "deal-general-info", templates.DealGeneralInfoArgs{Deal: deal})
}
