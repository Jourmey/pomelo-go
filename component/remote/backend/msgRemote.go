package backend

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"pomelo-go/cluster/clusterpb/proto"
	"pomelo-go/component"
	"pomelo-go/internal/codec"
)

const (
	namespace = "sys"
	service   = "msgRemote"
)

var _ = component.Component(&Component{})

type Component struct {
	component.Base

	serverType string
}

func NewComponent(serverType string) *Component {

	return &Component{
		Base:       component.Base{},
		serverType: serverType,
	}

}

func (c *Component) Routes() (router map[string]component.Handler) {

	route := fmt.Sprintf("%s.%s.%s.%s", namespace, c.serverType, service, "forwardMessage")
	router = map[string]component.Handler{
		route: c.forwardMessageHandler,
	}

	return router
}

//func (c *Component) AddRoutes(rs []component.Route) {
//
//	c.router = append(c.router, rs...)
//}

func (c *Component) forwardMessageHandler(ctx context.Context, in json.RawMessage) (out json.RawMessage) {

	session := proto.Session{}
	msg := proto.Message{}

	if err := codec.Decode(in, []interface{}{&msg, &session}); err != nil {
		return codec.Encode(err.Error(), nil)
	}

	res := map[string]interface{}{
		"a": "A",
		"b": "BBB",
	}

	return codec.Encode(errors.New("test panic").Error(), res)
	//return codec.Encode(nil, res)
}
