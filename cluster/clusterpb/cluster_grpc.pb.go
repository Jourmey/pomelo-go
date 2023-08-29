package clusterpb

import (
	"context"
	"pomelo-go/cluster/clusterpb/proto"
)

type MasterClient interface {
	// Register 向master注册服务信息
	Register(ctx context.Context, in *proto.RegisterRequest) (*proto.RegisterResponse, error)
	// Subscribe 订阅master中集群信息
	Subscribe(ctx context.Context, in *proto.SubscribeRequest) (*proto.SubscribeResponse, error)
	// Record 通知master启动完毕
	Record(ctx context.Context, in *proto.RecordRequest) (*proto.RecordResponse, error)
	// MonitorHandler 监听master中的集群变化
	MonitorHandler(ctx context.Context, in *proto.MonitorHandlerRequest) (*proto.MonitorHandlerResponse, error)
}

//type MasterServer interface {
//}

type MemberClient interface {
	// Request 发送Request rpc请求
	Request(ctx context.Context, in *proto.RequestRequest) (*proto.RequestResponse, error)
	// Notify 发送Notify rpc请求
	Notify(ctx context.Context, in *proto.NotifyRequest) (*proto.NotifyResponse, error)
}

type MemberServer interface {
}

type MasterAgent interface {
	MasterClient

	Connect() error
	Close() error
}

type MemberAgent interface {
	MemberClient

	Connect() error
	Close() error
}
