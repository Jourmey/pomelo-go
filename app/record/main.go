package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	component2 "pomelo-go/app/record/internal/component"
	"pomelo-go/app/record/internal/config"
	"pomelo-go/app/record/internal/svc"
	"pomelo-go/component"
	"pomelo-go/pomelo"
	"time"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {

	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)

	// 业务层 服务组件
	components := &component.Components{}
	components.Register(component2.NewRecord(ctx))

	pomelo.Listen(c.Listen, // 本地服务rpc地址
		pomelo.WithAdvertiseAddr(c.AdvertiseAddr),                                           // node服务对应的master地址
		pomelo.WithAdvertiseRetry(time.Duration(c.RetryInterval)*time.Second, c.RetryTimes), // 注册重试配置
		pomelo.WithServerId(c.Name),                                                         // 本机服务名称
		pomelo.WithServerInfo(c.ServerInfo),                                                 // 本机服务信息配置
		pomelo.WithToken(c.Token),                                                           // 与master通信token
		pomelo.WithComponents(components),
	)

}
