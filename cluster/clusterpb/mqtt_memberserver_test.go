package clusterpb

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"pomelo-go/cluster/clusterpb/proto"
	"testing"
)

type myMemberServer struct {
}

func (m *myMemberServer) RequestHandler(ctx context.Context, in *proto.RequestRequest) (*proto.RequestResponse, error) {
	logx.Info("myMemberServer RequestHandler")
	return &proto.RequestResponse{}, nil
}

func (m *myMemberServer) NotifyHandler(ctx context.Context, in *proto.NotifyRequest) (*proto.NotifyResponse, error) {
	logx.Info("myMemberServer NotifyHandler")
	return &proto.NotifyResponse{}, nil
}

func Test_NewMqttMasterServer(t *testing.T) {

	var (
		advertiseAddr = ":8081"
	)

	m := myMemberServer{}

	server := NewMqttMasterServer(&m)

	err := server.Listen(advertiseAddr)
	if err != nil {
		t.Fatal(err)
	}

	select {}
}
