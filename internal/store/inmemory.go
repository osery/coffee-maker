package store

import (
	"sync"

	"github.com/osery/coffee-maker/pkg/model"
)

func NewInMemory() Store {
	return &inMemory{
		coffees:       []*model.Coffee{},
		coffeesByName: map[string]*model.Coffee{},
	}
}

type inMemory struct {
	mu            sync.Mutex
	coffees       []*model.Coffee
	coffeesByName map[string]*model.Coffee
}

func (i *inMemory) ListCoffees() ([]*model.Coffee, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	result := make([]*model.Coffee, len(i.coffees))
	for i, coffee := range i.coffees {
		result[i] = coffee.DeepCopy()
	}

	return result, nil
}

func (i *inMemory) GetCoffeeByName(name string) (*model.Coffee, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	coffee, in := i.coffeesByName[name]
	if !in {
		return nil, ErrNotFound
	}

	return coffee.DeepCopy(), nil
}

func (i *inMemory) InsertCoffee(coffee *model.Coffee) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	_, in := i.coffeesByName[coffee.Name]
	if in {
		return ErrDuplicateRecord
	}
	i.coffees = append(i.coffees, coffee)
	i.coffeesByName[coffee.Name] = coffee
	return nil
}

func (i *inMemory) UpdateCoffeeStatus(name string, status model.CoffeeStatus) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	coffee, in := i.coffeesByName[name]
	if !in {
		return ErrNotFound
	}
	coffee.Status = status
	return nil
}
