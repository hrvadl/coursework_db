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

// TODO: my services are useless they should not copy repositories
// all logic in controllers should have been located in services
func NewTransaction(
	us services.User,
	ds services.Deal,
	is services.Inventory,
	tr *templates.Resolver,
) *Transaction {
	return &Transaction{
		ss: us,
		ds: ds,
		is: is,
		tr: tr,
	}
}

type Transaction struct {
	ss services.User
	ds services.Deal
	is services.Inventory
	tr *templates.Resolver
}

type userStrategy interface {
	Patch(u *models.User) (*models.User, error)
}

func (t *Transaction) HandleMakeTransaction(w http.ResponseWriter, r *http.Request) {
	ctx := middleware.Must(
		middleware.GetUserCtx(r.Context()),
	)

	parts := strings.Split(r.URL.Path, "/")
	dealID, err := strconv.ParseInt(parts[len(parts)-1], 10, 64)

	if err != nil {
		t.tr.Execute(w, "generic-error.html", templates.GenericErrorArgs{})
		return
	}

	r.ParseForm()
	amount, err := strconv.ParseInt(r.FormValue("amount"), 10, 64)

	if err != nil {
		t.tr.Execute(w, "toast", templates.ToastArgs{Error: "invalid amount"})
		return
	}

	err = t.ds.MakeTransaction(int(ctx.ID), int(dealID), int(amount))

	if err != nil {
		t.tr.Execute(w, "toast", templates.ToastArgs{Error: err.Error()})
		return
	}

}

func (t *Transaction) handleBuying(w http.ResponseWriter, r *http.Request) {

}

func (t *Transaction) handleSelling(w http.ResponseWriter, r *http.Request) {}
