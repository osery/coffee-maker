package store

import (
	"errors"

	"github.com/osery/coffee-maker/pkg/model"
)

var (
	ErrNotFound        = errors.New("record not found")
	ErrDuplicateRecord = errors.New("duplicate record")
)

type Store interface {
	ListCoffees() ([]*model.Coffee, error)
	GetCoffeeByName(name string) (*model.Coffee, error)
	InsertCoffee(coffee *model.Coffee) error
	UpdateCoffeeStatus(name string, status model.CoffeeStatus) error
}
