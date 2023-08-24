package cluster

import (
	"fmt"
	"pomelo-go/cluster/clusterpb/proto"
	"pomelo-go/component"
	"testing"
	"time"
)

type MyComponent struct {
}

func (m MyComponent) Init() {
	fmt.Println("MyComponent.Init")
}

func (m MyComponent) AfterInit() {
	fmt.Println("MyComponent.AfterInit")
}

func (m MyComponent) BeforeShutdown() {
	fmt.Println("MyComponent.BeforeShutdown")
}

func (m MyComponent) Shutdown() {
	fmt.Println("MyComponent.Shutdown")
}

func TestNode_Startup(t *testing.T) {

	c := &component.Components{}
	c.Register(&MyComponent{})

	serverid := "cluster-server-connector-996"

	opt := Options{
		IsMaster:      false,
		ServerId:      serverid,
		AdvertiseAddr: "localhost:3005",
		ServerInfo: proto.ClusterServerInfo{
			Id:         serverid,
			Type:       proto.Type_Monitor,
			ServerType: proto.ServerType_Connector,
			Pid:        99,
			Info: map[string]interface{}{
				"env":           "local",
				"host":          "127.0.0.1",
				"port":          4061,
				"clientPort":    3061,
				"wssPort":       80,
				"frontend":      "false",
				"channelType":   2,
				"cloudType":     1,
				"clusterCount":  1,
				"restart-force": "true",
				"recover":       "true",
				"serverType":    proto.ServerType_Connector,
				"id":            serverid,
			},
		},
		RetryInterval: 5 * time.Second,
		RetryTimes:    60,
		Token:         "agarxhqb98rpajloaxn34ga8xrunpagkjwlaw3ruxnpaagl29w4rxn",
		Components:    c,
	}

	n := &Node{
		Options:     opt,
		ServiceAddr: "127.0.0.1:4450",
	}

	err := n.Startup()
	if err != nil {
		t.Fatal(err)
	}

	select {}
}
