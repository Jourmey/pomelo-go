package pomelo

import (
	"pomelo-go/cluster"
	"pomelo-go/cluster/clusterpb/proto"
	"pomelo-go/component"
	"time"
)

type Option func(*cluster.Options)

func WithAdvertiseAddr(addr string) Option {
	return func(opt *cluster.Options) {
		opt.AdvertiseAddr = addr
	}
}

func WithAdvertiseRetry(retryInterval time.Duration, retryTimes int) Option {
	return func(opt *cluster.Options) {
		opt.RetryInterval = retryInterval
		opt.RetryTimes = retryTimes
	}
}

func WithMaster() Option {
	return func(opt *cluster.Options) {
		opt.IsMaster = true
	}
}

func WithComponents(components *component.Components) Option {
	return func(opt *cluster.Options) {
		opt.Components = components
	}
}

func WithToken(token string) Option {
	return func(opt *cluster.Options) {
		opt.Token = token
	}
}

func WithServerId(serverId string) Option {
	return func(opt *cluster.Options) {
		opt.ServerId = serverId
	}
}

func WithServerInfo(ServerInfo proto.ClusterServerInfo) Option {
	return func(opt *cluster.Options) {
		opt.ServerInfo = ServerInfo
	}
}
