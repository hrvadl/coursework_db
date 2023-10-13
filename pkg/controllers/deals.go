package controllers

import (
	"net/http"

	"github.com/hrvadl/coursework_db/pkg/services"
	"github.com/hrvadl/coursework_db/pkg/templates"
)

func NewDeal(ds services.Deal, t *templates.Resolver) *Deal {
	return &Deal{
		ds: ds,
		t:  t,
	}
}

type Deal struct {
	ds services.Deal
	t  *templates.Resolver
}

func (d *Deal) ServeDealsPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, d.t.ResolveHTML("deals.html"))
}
