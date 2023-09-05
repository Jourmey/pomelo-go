package types

type (
	// MsgRecoverTemporaryRequest 互动恢复路由 压缩
	MsgRecoverTemporaryRequest struct {
		LiveId     string `json:"liveId"`
		LecturerId string `json:"lecturerId"`
		TutorId    string `json:"tutorId"`
		StuId      string `json:"stuId"`
	}

	MsgRecoverTemporaryResponse struct {
		Str   string
		Slice []string
		Map   map[string]interface{}
	}
)
type (
	// MsgRecoverClassroomRequest 互动恢复路由 压缩 座位席恢复  可以通过课堂数据导入  分全量和增量
	MsgRecoverClassroomRequest  struct{}
	MsgRecoverClassroomResponse struct{}
)
type (
	// MsgRecoverDataInClassRequest 互动恢复路由 压缩 同屏恢复
	MsgRecoverDataInClassRequest  struct{}
	MsgRecoverDataInClassResponse struct{}
)
type (
	// MsgRecoverStreamRequest rtc频道 信令消息恢复路由 4分钟过期
	MsgRecoverStreamRequest  struct{}
	MsgRecoverStreamResponse struct{}
)
