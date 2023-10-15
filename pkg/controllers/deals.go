package controllers

import (
	"net/http"

	"github.com/hrvadl/coursework_db/pkg/middleware"
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
