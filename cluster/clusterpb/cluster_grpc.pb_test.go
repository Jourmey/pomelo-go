package clusterpb

import (
	"context"
	"log"
	"pomelo-go/cluster/clusterpb/proto"
	"testing"
	"time"
)

var (
	host = "localhost"
	port = 3005

	serverId = "cluster-server-connector-998"

	request = &proto.RegisterRequest{
		ServerInfo: proto.ClusterServerInfo{
			"id":         serverId,
			"type":       proto.Type_Monitor,
			"serverType": proto.ServerType_Connector,
			"pid":        99,
			"info": map[string]interface{}{
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
				"serverType":    proto.ServerType_Connector,
				"id":            serverId,
			},
		},
		Token: "agarxhqb98rpajloaxn34ga8xrunpagkjwlaw3ruxnpaagl29w4rxn",
	}

	client MasterClient
)

func init() {
	c := NewMasterClient(host, port)

	for {
		err := c.Connect()
		if err == nil {
			break
		}

		time.Sleep(5 * time.Second)
		log.Println("try connect again")
	}

	client = c
}

func Test_MqttMasterClient_All(t *testing.T) {

	Test_MqttMasterClient_Register(t)

	Test_MqttMasterClient_Subscribe(t)

	Test_MqttMasterClient_Record(t)
}

func Test_MqttMasterClient_Register(t *testing.T) {

	res, err := client.Register(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(res)
}

func Test_MqttMasterClient_Subscribe(t *testing.T) {

	res, err := client.Subscribe(context.Background(), &proto.SubscribeRequest{
		Id: serverId,
	})

	if err != nil {
		t.Fatal(err)
	}

	t.Log(res)
}

func Test_MqttMasterClient_Record(t *testing.T) {

	res, err := client.Record(context.Background(), &proto.RecordRequest{
		Id: serverId,
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Log(res)
}
