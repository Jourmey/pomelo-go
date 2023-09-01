package config

import "pomelo-go/cluster/clusterpb/proto"

type Config struct {
	Name          string                  // 服务名称
	Listen        string                  // 本地服务rpc地址
	AdvertiseAddr string                  // node服务对应的master地址
	ServerInfo    proto.ClusterServerInfo // 本机服务信息  TODO 自动获取
	RetryInterval int                     // 注册重试间隔（单位 s）
	RetryTimes    int                     // 注册重试次数
	Token         string                  // 与master通信token
}
