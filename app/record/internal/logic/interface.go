package logic

import "pomelo-go/app/record/internal/types"

type Record interface {
	MsgRecoverTemporary(req types.MsgRecoverTemporaryRequest) (res types.MsgRecoverTemporaryResponse, err error)       // 互动恢复路由 压缩
	MsgRecoverClassroom(req types.MsgRecoverClassroomRequest) (res types.MsgRecoverClassroomResponse, err error)       // 互动恢复路由 压缩 座位席恢复  可以通过课堂数据导入  分全量和增量
	MsgRecoverDataInClass(req types.MsgRecoverDataInClassRequest) (res types.MsgRecoverDataInClassResponse, err error) // 互动恢复路由 压缩 同屏恢复
	MsgRecoverStream(req types.MsgRecoverStreamRequest) (res types.MsgRecoverStreamResponse, err error)                // rtc频道 信令消息恢复路由 4分钟过期
}
