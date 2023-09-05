package clusterpb

import (
	"context"
	"encoding/json"
	"log"
	"pomelo-go/cluster/clusterpb/proto"
	"testing"
	"time"
)

var (
	memberClient MemberClientAgent
)

func InitMqttMemberClient() {
	var (
		advertiseAddr = "127.0.0.1:10061"
	)

	c := NewMqttMemberClient(advertiseAddr)

	for {
		err := c.Connect()
		if err == nil {
			break
		}

		time.Sleep(5 * time.Second)
		log.Println("try connect again")
	}

	memberClient = c
}

func Test_MqttMemberClient_Request(t *testing.T) {

	InitMqttMemberClient()

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

	res, err := memberClient.Request(context.Background(), proto.RequestRequest{
		Namespace:  "user",
		ServerType: "chat",
		Service:    "chatRemote",
		Method:     "add",
		Args:       []json.RawMessage{body},
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(res)
}
