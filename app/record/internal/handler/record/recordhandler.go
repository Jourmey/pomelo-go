package record

import (
	"context"
	"encoding/json"
	"pomelo-go/app/record/internal/logic/record"
	"pomelo-go/app/record/internal/svc"
	"pomelo-go/app/record/internal/types"
	"pomelo-go/cluster/clusterpb/proto"
	"pomelo-go/component/remote/backend"
)

func MsgRecoverTemporaryHandler(serverCtx *svc.ServiceContext) backend.ForwardMessageHandler {
	return func(ctx context.Context, session proto.Session, message proto.Message) (interface{}, error) {
		var req types.MsgRecoverTemporaryRequest
		if err := json.Unmarshal(message.Body, &req); err != nil {
			return nil, err
		}

		l := record.NewRecordLogic(ctx, serverCtx)
		resp, err := l.MsgRecoverTemporary(req)
		return resp, err
	}
}

func MsgRecoverClassroomHandler(serverCtx *svc.ServiceContext) backend.ForwardMessageHandler {
	return func(ctx context.Context, session proto.Session, message proto.Message) (interface{}, error) {
		var req types.MsgRecoverClassroomRequest
		if err := json.Unmarshal(message.Body, &req); err != nil {
			return nil, err
		}

		l := record.NewRecordLogic(ctx, serverCtx)
		resp, err := l.MsgRecoverClassroom(req)
		return resp, err
	}
}
