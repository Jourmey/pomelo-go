package pomelo

import (
	"context"
	"encoding/json"
	"pomelo-go/cluster/clusterpb/proto"
	"testing"
	"time"
)

func Test_Listen(t *testing.T) {

	time.AfterFunc(5*time.Second, func() {
		Test_RemoteProcess(t)
	})

	serverid := "cluster-server-record-996"

	Listen("127.0.0.1:8081", // 本机rpc服务addr
		WithServerId(serverid),
		WithServerInfo(proto.ClusterServerInfo{
			"type":          proto.Type_Monitor,
			"pid":           99,
			"env":           "local",
			"host":          "127.0.0.1",
			"port":          8081,
			"channelType":   2,
			"cloudType":     1,
			"clusterCount":  1,
			"restart-force": "true",
			"serverType":    proto.ServerType_Recover,
			"id":            serverid,
		}),
		WithToken("agarxhqb98rpajloaxn34ga8xrunpagkjwlaw3ruxnpaagl29w4rxn"),
		WithAdvertiseAddr("localhost:3005"), // 集群master服务addr
	)
}

func Test_RemoteProcess(t *testing.T) {

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

	res, err := RemoteProcess(context.Background(), &proto.RequestRequest{
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

}
