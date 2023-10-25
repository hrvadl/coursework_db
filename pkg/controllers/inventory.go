package controllers

import (
	"net/http"
	"strconv"

	"github.com/hrvadl/coursework_db/pkg/middleware"
	"github.com/hrvadl/coursework_db/pkg/services"
	"github.com/hrvadl/coursework_db/pkg/templates"
)

func NewInventory(
	is services.Inventory,
	us services.User,
	tr *templates.Resolver,
) *Inventory {
	return &Inventory{
		is: is,
		tr: tr,
		us: us,
	}
}

type Inventory struct {
	is services.Inventory
	us services.User
	tr *templates.Resolver
}

func (i *Inventory) HandlePatch(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.Must(
		middleware.GetUserCtx(r.Context()),
	)

	r.ParseForm()
	securityID, err := strconv.ParseInt(r.FormValue("securityID"), 10, 64)

	if err != nil {
		i.tr.Execute(w, "generic-error.html", templates.GenericErrorArgs{})
		return
	}

	amount, err := strconv.ParseInt(r.FormValue("amount"), 10, 64)

	if err != nil {
		i.tr.Execute(w, "generic-error.html", templates.GenericErrorArgs{})
		return
	}

	withdraw, _ := strconv.ParseBool(r.FormValue("withdraw"))
	switch withdraw {
	case true:
		_, err = i.is.Withdraw(ctx.ID, uint(securityID), uint(amount))
	case false:
		_, err = i.is.Add(ctx.ID, uint(securityID), uint(amount))
	}

	if err != nil {
		i.tr.Execute(w, "toast", templates.ToastArgs{Error: "Something went wrong"})
		return
	}

	w.Header().Set("HX-Trigger", "refresh-inventory-info")
	w.WriteHeader(http.StatusCreated)
}

func (i *Inventory) HandleGetInventoryInfo(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.Must(
		middleware.GetUserCtx(r.Context()),
	)

	user, err := i.us.GetByID(int(ctx.ID))

	if err != nil {
		i.tr.Execute(w, "generic-error.html", templates.GenericErrorArgs{})
		return
	}

	i.tr.Execute(w, "inventory", templates.InventoryArgs{
		User: user,
	})
}
