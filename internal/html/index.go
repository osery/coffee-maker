package html

import (
	"html/template"
	"net/http"

	"github.com/osery/coffee-maker/internal/store"
	"github.com/osery/coffee-maker/internal/util"
)

func NewHTMLHandler(store store.Store) *handler {
	return &handler{
		store: store,
	}
}

type handler struct {
	store store.Store
}

func (h *handler) IndexPage(w http.ResponseWriter, r *http.Request) {
	coffees, err := h.store.ListCoffees()
	if err != nil {
		util.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	t, err := template.ParseFiles("static/index.html.tmpl")
	if err != nil {
		util.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, coffees)
	if err != nil {
		util.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
}
