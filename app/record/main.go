package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
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

	components := component.NewComponents()
	components.AddRoutes(
		[]component.Route{
			{"sys.recover.msgRemote.forwardMessage", forwardMessageHandler(ctx)},
		})

	pomelo.Listen(c.Listen, // 本地服务rpc地址
		pomelo.WithAdvertiseAddr(c.AdvertiseAddr),                                           // node服务对应的master地址
		pomelo.WithAdvertiseRetry(time.Duration(c.RetryInterval)*time.Second, c.RetryTimes), // 注册重试配置
		pomelo.WithServerId(c.Name),                                                         // 本机服务名称
		pomelo.WithServerInfo(c.ServerInfo),                                                 // 本机服务信息配置
		pomelo.WithToken(c.Token),                                                           // 与master通信token
		pomelo.WithComponents(components),                                                   // 业务层 服务组件
	)

}

func forwardMessageHandler(svcCtx *svc.ServiceContext) component.Handler {
	return func(ctx context.Context, in []json.RawMessage) (out []json.RawMessage) {

		l := NewOperator(ctx, svcCtx)
		res, err := l.forwardMessage(in)
		e, _ := json.Marshal(err)
		r, _ := json.Marshal(res)

		return []json.RawMessage{e, r}

	}
}

type MsgRemote struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func (r *MsgRemote) forwardMessage(in []json.RawMessage) (interface{}, error) {

	r.Logger.Infof("msgRemote forwardMessage in:%s", in)

	res := map[string]interface{}{
		"a": "A",
		"b": "BBB",
	}

	return res, nil
}

func NewOperator(ctx context.Context, svcCtx *svc.ServiceContext) *MsgRemote {
	logger := logx.WithContext(ctx)
	return &MsgRemote{
		Logger: logger,
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}
