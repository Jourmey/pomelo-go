package component

import (
	"context"
	"encoding/json"
)

type Handler func(ctx context.Context, in json.RawMessage) (out json.RawMessage)

type Component interface {
	Init()
	AfterInit()
	BeforeShutdown()
	Shutdown()
	Routes() (router map[string]Handler)
}

// A Route is a http route.
type Route struct {
	Method  string
	Handler Handler
}
