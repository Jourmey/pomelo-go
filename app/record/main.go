package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"pomelo-go/app/record/internal/config"
	"pomelo-go/app/record/internal/handler"
	"pomelo-go/app/record/internal/svc"
	"pomelo-go/cluster"
	"pomelo-go/component/remote/backend"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)

	forwardMessage := backend.NewComponent("recover")
	server := cluster.MustNewServer(c.Config, cluster.WithComponent(forwardMessage))
	defer server.Stop()

	handler.RegisterHandlers(forwardMessage, ctx) // 路由handler注册

	fmt.Printf("Starting server at %s...\n", c.AdvertiseAddr)
	server.Start()
}
