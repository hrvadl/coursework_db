package controllers

import (
	"net/http"
	"strconv"

	"github.com/hrvadl/coursework_db/pkg/middleware"
	"github.com/hrvadl/coursework_db/pkg/models"
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

	if amount < 1 {
		i.tr.Execute(w, "toast", templates.ToastArgs{Error: "Invalid money amount"})
		return
	}

	inventory, _ := i.is.GetUserInventoryBySecurityID(int(ctx.ID), int(securityID))

	if inventory == nil {
		inventory = &models.InventoryItem{}
	}

	inventory.SecurityID = uint(securityID)
	inventory.OwnerID = ctx.ID

	withdraw, _ := strconv.ParseBool(r.FormValue("withdraw"))
	switch withdraw {
	case true:
		inventory.Amount -= uint(amount)
	case false:
		inventory.Amount += uint(amount)
	}

	if _, err := i.is.CreateOrUpdate(inventory); err != nil {
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
