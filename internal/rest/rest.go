package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/osery/coffee-maker/internal/brew"
	"github.com/osery/coffee-maker/internal/store"
	"github.com/osery/coffee-maker/internal/util"
	"github.com/osery/coffee-maker/pkg/model"
)

func NewRESTHandler(store store.Store, brewer brew.Brewer) *handler {
	return &handler{
		store:  store,
		brewer: brewer,
	}
}

type handler struct {
	store  store.Store
	brewer brew.Brewer
}

func (h *handler) ListCoffees(w http.ResponseWriter, _ *http.Request) {
	coffees, err := h.store.ListCoffees()
	if err != nil {
		util.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	util.SuccessResponse(w, coffees)
}

func (h *handler) GetCoffee(w http.ResponseWriter, r *http.Request) {
	name, err := util.GetStringVar("name", r)
	if err != nil {
		util.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	coffee, err := h.store.GetCoffeeByName(name)
	if errors.Is(err, store.ErrNotFound) {
		util.ErrorResponse(w, err, http.StatusNotFound)
		return
	}
	if err != nil {
		util.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	util.SuccessResponse(w, coffee)
}

func (h *handler) PostCoffee(w http.ResponseWriter, r *http.Request) {
	coffee := new(model.Coffee)
	err := json.NewDecoder(r.Body).Decode(coffee)
	if err != nil {
		util.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}
	err = checkCoffee(coffee)
	if err != nil {
		util.ErrorResponse(w, err, http.StatusBadRequest)
		return
	}

	coffee.Status = model.Queued

	err = h.store.InsertCoffee(coffee)
	if err != nil {
		util.ErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	err = h.brewer.BrewCoffee(coffee)
	if err != nil {
		util.ErrorResponse(w, err, http.StatusInternalServerError)
		err = h.store.UpdateCoffeeStatus(coffee.Name, model.Failed)
		if err != nil {
			zap.L().Error("Failed marking coffee as failed", zap.String("name", coffee.Name), zap.Error(err))
		}
	}
}

func checkCoffee(coffee *model.Coffee) error {
	if coffee.Name == "" {
		return fmt.Errorf("empty record name")
	}
	if !model.IsValidCoffeeType(coffee.Type) {
		return fmt.Errorf("invalid coffee type %s", coffee.Type)
	}
	return nil
}
