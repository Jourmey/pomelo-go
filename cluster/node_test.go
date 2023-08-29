package cluster

import (
	"context"
	"encoding/json"
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

	serverid := "cluster-server-chat-996"

	opt := Options{
		IsMaster:      false,
		ServerId:      serverid,
		AdvertiseAddr: "localhost:3005",
		ServerInfo: proto.ClusterServerInfo{
			"type":          proto.Type_Monitor,
			"pid":           99,
			"env":           "local",
			"host":          "127.0.0.1",
			"port":          4061,
			"clientPort":    3061,
			"wssPort":       80,
			"frontend":      "true",
			"channelType":   2,
			"cloudType":     1,
			"clusterCount":  1,
			"restart-force": "true",
			"recover":       "true",
			"serverType":    proto.ServerType_Chat,
			"id":            serverid,
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

	args := []interface{}{
		"stu1*kick_testsss",
		"cluster-server-connector-0",
		"kick_testsss",
		true,
		2,
		1,
		0,
		"123",
		"abc",
		"2.9.8.7",
		map[string]interface{}{
			"uniqId":    "231FF2BB-BA09-598D-9EB6-3B0299D292E7ssss",
			"rid":       "kick_testsss",
			"rtype":     2,
			"role":      1,
			"ulevel":    0,
			"uname":     "123",
			"classid":   "abc",
			"clientVer": "2.9.8.7",
			"userVer":   "1.0",
			"liveType":  "COMBINE_SMALL_CLASS_MODE"},
		"0"}

	body, err := json.Marshal(args)
	if err != nil {
		t.Fatal(err)
	}

	res, err := n.RemoteProcess(context.Background(), &proto.RequestRequest{
		Namespace:  "user",
		ServerType: "chat",
		Service:    "chatRemote",
		Method:     "add",
		Args:       body,
	})

	if err != nil {
		t.Fatal(err)
	}

	t.Log(res)

	select {}
}
