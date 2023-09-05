package handler

import (
	"context"
	"errors"
	record2 "pomelo-go/app/record/internal/handler/record"
	"pomelo-go/app/record/internal/svc"
	"pomelo-go/cluster/clusterpb/proto"
	"pomelo-go/component/remote/backend"
)

func RegisterHandlers(server *backend.Component, serverCtx *svc.ServiceContext) {

	server.AddRoutes([]backend.Route{
		{"recover.recoverHandler.msgRecoverTemporary", record2.MsgRecoverTemporaryHandler(serverCtx)}, // 互动恢复路由 压缩
		{"recover.recoverHandler.msgRecoverClassroom", record2.MsgRecoverClassroomHandler(serverCtx)}, // 互动恢复路由 压缩 座位席恢复  可以通过课堂数据导入  分全量和增量
		{"recover.recoverHandler.msgRecoverDataInClass", NilHandler(serverCtx)},                       // 互动恢复路由 压缩 同屏恢复
		{"recover.recoverHandler.msgRecoverStream", NilHandler(serverCtx)},                            // rtc频道 信令消息恢复路由 4分钟过期
	})
}

func NilHandler(ctx *svc.ServiceContext) backend.ForwardMessageHandler {
	return func(ctx context.Context, session proto.Session, message proto.Message) (interface{}, error) {
		return nil, errors.New("nil function")
	}
}
