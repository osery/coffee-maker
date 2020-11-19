package model

import (
	"github.com/jinzhu/copier"
)

type CoffeeType string

const (
	Americano CoffeeType = "americano"
	Espresso  CoffeeType = "espresso"
	Latte     CoffeeType = "latte"
)

type CoffeeStatus string

var (
	supportedCoffeeTypes = map[CoffeeType]bool{
		Americano: true,
		Espresso:  true,
		Latte:     true,
	}
)

const (
	Queued  CoffeeStatus = "queued"
	Brewing CoffeeStatus = "brewing"
	Done    CoffeeStatus = "done"
	Failed  CoffeeStatus = "failed"
)

type Coffee struct {
	Name       string
	Type       CoffeeType
	ExtraSugar bool `json:"extraSugar"`
	Status     CoffeeStatus
}

func IsValidCoffeeType(t CoffeeType) bool {
	_, in := supportedCoffeeTypes[t]
	return in
}

func (c *Coffee) DeepCopy() *Coffee {
	dup := &Coffee{}
	copier.Copy(dup, c)
	return dup
}
