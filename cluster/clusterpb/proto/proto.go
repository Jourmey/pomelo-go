package proto

type MonitorAction string

const (
	Type_Monitor = "monitor"

	ServerType_Connector = "connector"

	MonitorAction_addServer     MonitorAction = "addServer"
	MonitorAction_removeServer  MonitorAction = "removeServer"
	MonitorAction_replaceServer MonitorAction = "replaceServer"
	MonitorAction_startOve      MonitorAction = "startOver"
)

const (
	BEFORE_FILTER        = "__befores__"
	AFTER_FILTER         = "__afters__"
	GLOBAL_BEFORE_FILTER = "__globalBefores__"
	GLOBAL_AFTER_FILTER  = "__globalAfters__"
	ROUTE                = "__routes__"
	BEFORE_STOP_HOOK     = "__beforeStopHook__"
	MODULE               = "__modules__"
	SERVER_MAP           = "__serverMap__"
	RPC_BEFORE_FILTER    = "__rpcBefores__"
	RPC_AFTER_FILTER     = "__rpcAfters__"
	MASTER_WATCHER       = "__masterwatcher__"
	MONITOR_WATCHER      = "__monitorwatcher__"
)

// ClusterServerInfo 集群服务信息
type ClusterServerInfo map[string]interface{}

//type ClusterServerInfo struct {
//	Id         string                 `json:"id"`
//	Type       string                 `json:"type"`
//	ServerType string                 `json:"serverType"`
//	Pid        int                    `json:"pid"`
//	Info       map[string]interface{} `json:"info"`
//}

// Register 向master注册服务信息
type (
	RegisterRequest struct {
		ServerInfo ClusterServerInfo

		Token string `json:"token"`
	}

	RegisterResponse struct{}
)

// Subscribe 订阅master中集群信息
type (
	SubscribeRequest struct {
		Id string `json:"id"`
	}

	SubscribeResponse map[string]ClusterServerInfo // 集群内其他服务信息
)

// Record 通知master启动完毕
type (
	RecordRequest struct {
		Id string `json:"id"`
	}

	RecordResponse struct{}
)

// MonitorHandler 监听master中的集群变化
type (
	MonitorHandlerRequest struct {
		CallBackHandler func(action MonitorAction, serverInfos []ClusterServerInfo)
	}

	MonitorHandlerResponse struct{}
)
