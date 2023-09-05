package component

import (
	"context"
	"encoding/json"
)

type Handler func(ctx context.Context, in []json.RawMessage) (out []json.RawMessage)

type Components struct {
	Router map[string]Handler
}

// A Route is a http route.
type Route struct {
	Router  string
	Handler Handler
}

func NewComponents() *Components {
	return &Components{
		Router: make(map[string]Handler, 0),
	}
}

// AddRoutes add given routes into the Server.
func (cs *Components) AddRoutes(rs []Route) {

	for i := 0; i < len(rs); i++ {

		cs.Router[rs[i].Router] = rs[i].Handler
	}
}
