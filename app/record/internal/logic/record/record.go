package record

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"pomelo-go/app/record/internal/logic"
	"pomelo-go/app/record/internal/svc"
	"pomelo-go/app/record/internal/types"
)

type recordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (r *recordLogic) MsgRecoverTemporary(req types.MsgRecoverTemporaryRequest) (res types.MsgRecoverTemporaryResponse, err error) {

	return types.MsgRecoverTemporaryResponse{
		Str:   "AAA",
		Slice: []string{"BBB", "CCC"},
		Map: map[string]interface{}{
			"DDD": "EEEE",
		},
	}, nil

}

func (r *recordLogic) MsgRecoverClassroom(req types.MsgRecoverClassroomRequest) (res types.MsgRecoverClassroomResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *recordLogic) MsgRecoverDataInClass(req types.MsgRecoverDataInClassRequest) (res types.MsgRecoverDataInClassResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *recordLogic) MsgRecoverStream(req types.MsgRecoverStreamRequest) (res types.MsgRecoverStreamResponse, err error) {
	//TODO implement me
	panic("implement me")
}

func NewRecordLogic(ctx context.Context, svcCtx *svc.ServiceContext) logic.Record {
	logger := logx.WithContext(ctx)
	return &recordLogic{
		Logger: logger,
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
