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

//type MemberClient interface {
//}

//type MemberServer interface {
//	HandleRequest(context.Context, *RequestMessage) (*MemberHandleResponse, error)
//	HandleNotify(context.Context, *NotifyMessage) (*MemberHandleResponse, error)
//	HandlePush(context.Context, *PushMessage) (*MemberHandleResponse, error)
//	HandleResponse(context.Context, *ResponseMessage) (*MemberHandleResponse, error)
//	NewMember(context.Context, *NewMemberRequest) (*NewMemberResponse, error)
//	DelMember(context.Context, *DelMemberRequest) (*DelMemberResponse, error)
//	SessionClosed(context.Context, *SessionClosedRequest) (*SessionClosedResponse, error)
//	CloseSession(context.Context, *CloseSessionRequest) (*CloseSessionResponse, error)
//}
