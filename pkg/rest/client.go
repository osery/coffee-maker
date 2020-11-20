package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"

	"github.com/osery/coffee-maker/pkg/model"
)

type Client interface {
	Create(ctx context.Context, name string, t model.CoffeeType, extraSugar bool) (*model.Coffee, error)
	GetByName(ctx context.Context, name string) (*model.Coffee, error)
}

func NewClient(baseURL string) Client {
	return &client{
		resty:   resty.New(),
		baseURL: baseURL,
	}
}

type client struct {
	resty   *resty.Client
	baseURL string
}

func (c *client) Create(ctx context.Context, name string, t model.CoffeeType, extraSugar bool) (*model.Coffee, error) {
	var coffee model.Coffee
	resp, err := c.resty.R().
		SetContext(ctx).
		SetBody(model.Coffee{
			Name:       name,
			Type:       t,
			ExtraSugar: extraSugar,
		}).
		SetResult(&coffee).
		Post(fmt.Sprintf("%s/coffees", c.baseURL))
	if err != nil {
		return nil, fmt.Errorf("creating coffee (%s, %s, %t): %v", name, t, extraSugar, err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("creating coffee (%s, %s, %t) failed with status: %s", name, t, extraSugar, resp.Status())
	}
	return &coffee, nil
}

func (c *client) GetByName(ctx context.Context, name string) (*model.Coffee, error) {
	var coffee model.Coffee
	resp, err := c.resty.R().
		SetContext(ctx).
		SetResult(&coffee).
		Get(fmt.Sprintf("%s/coffees/%s", c.baseURL, name))
	if err != nil {
		return nil, fmt.Errorf("getting coffee (%s): %v", name, err)
	}
	if resp.StatusCode() == http.StatusNotFound {
		return nil, nil
	}
	if resp.IsError() {
		return nil, fmt.Errorf("getting coffee (%s) failed with status: %s", name, resp.Status())
	}
	return &coffee, nil
}
