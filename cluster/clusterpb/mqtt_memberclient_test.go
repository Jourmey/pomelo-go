package clusterpb

import (
	"context"
	"log"
	"pomelo-go/cluster/clusterpb/proto"
	"testing"
	"time"
)

var (
	memberClient MemberClient
)

func InitMqttMemberClient() {
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

	res, err := memberClient.Request(context.Background(), &proto.RequestRequest{})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(res)

}
