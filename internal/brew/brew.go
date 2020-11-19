package brew

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/osery/coffee-maker/internal/store"
	"github.com/osery/coffee-maker/pkg/model"
)

const QueueSize = 100

type Brewer interface {
	BrewCoffee(coffee *model.Coffee) error
}

func NewBrewer(store store.Store) Brewer {
	b := &brewer{
		store: store,
		queue: make(chan string, QueueSize),
	}
	go b.start()
	return b
}

type brewer struct {
	store store.Store
	queue chan string
}

func (b *brewer) BrewCoffee(coffee *model.Coffee) error {
	select {
	case b.queue <- coffee.Name:
		return nil
	default:
		return fmt.Errorf("coffee queue is full")
	}
}

func (b *brewer) start() {
	for {
		name := <-b.queue

		// Simulate brewing.
		b.updateStatus(name, model.Brewing)
		time.Sleep(time.Second * 30)
		b.updateStatus(name, model.Done)
	}
}

func (b *brewer) updateStatus(name string, status model.CoffeeStatus) {
	err := b.store.UpdateCoffeeStatus(name, status)
	if err != nil {
		zap.L().Error(
			"Failed changing coffee status",
			zap.String("name", name),
			zap.String("status", string(status)),
			zap.Error(err))
	}
}
