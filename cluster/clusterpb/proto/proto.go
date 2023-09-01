package proto

import "encoding/json"

type MonitorAction string

const (
	Type_Monitor = "monitor"

	ServerType_Connector = "connector"
	ServerType_Chat      = "chat"
	ServerType_Recover   = "recover"

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

// Register 向master注册服务信息
type (
	RegisterRequest struct {
		ServerInfo ClusterServerInfo

		Token string
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

type RpcMessage struct {
	Namespace  string          `json:"namespace"`
	ServerType string          `json:"serverType"`
	Service    string          `json:"service"`
	Method     string          `json:"method"`
	Args       json.RawMessage `json:"args"`
}

// Request 发送Request rpc请求
type (
	RequestRequest RpcMessage

	RequestResponse interface{}
)

// Notify 发送Notify rpc请求
type (
	NotifyRequest RpcMessage

	NotifyResponse struct{}
)
