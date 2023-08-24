package cluster

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"pomelo-go/cluster/clusterpb"
	"pomelo-go/cluster/clusterpb/proto"
	"pomelo-go/component"
	"time"
)

// Options contains some configurations for current node
type Options struct {
	IsMaster      bool
	ServerId      string                  // node服务id名称
	AdvertiseAddr string                  // node服务对应的master地址
	ServerInfo    proto.ClusterServerInfo // node 服务信息用于向master注册
	RetryInterval time.Duration           // master 重试间隔 default 5*time.sec
	RetryTimes    int                     // master 重试间隔 default 60次
	Token         string                  // master 通信token

	Components *component.Components
	//RemoteServiceRoute CustomerRemoteServiceRoute
}

type Node struct {
	Options            // current node options
	ServiceAddr string // current server service address (RPC)

	handler      *LocalHandler          // 处理本地或远程Handler调用
	masterClient clusterpb.MasterClient // 与master通信的客户端 对应pomelo的monitor
	//server  *grpcServer            // rpc服务端

	//sessions map[int64]*session.Session
}

func (n *Node) Startup() error {
	if n.ServiceAddr == "" {
		return errors.New("service address cannot be empty in master node")
	}
	n.handler = NewHandler(n)
	components := n.Components.List()
	for _, c := range components {
		err := n.handler.register(c.Comp, c.Opts)
		if err != nil {
			return err
		}
	}

	if err := n.initMasterClient(); err != nil {
		return err
	}

	_, err := n.masterClient.Register(context.Background(), &proto.RegisterRequest{
		ServerInfo: n.ServerInfo,
		Token:      n.Token,
	})
	if err != nil {
		return err
	}

	// 获取注册信息
	subscribeResponse, err := n.masterClient.Subscribe(context.Background(), &proto.SubscribeRequest{
		Id: n.ServerId,
	})

	// 初始化handler
	n.handler.initRemoteService(*subscribeResponse)

	// Initialize all components
	for _, c := range components {
		c.Comp.Init()
	}
	for _, c := range components {
		c.Comp.AfterInit()
	}

	return nil
}

func (n *Node) Handler() *LocalHandler {
	return n.handler
}

func (n *Node) Shutdown() {

}

func (n *Node) initMasterClient() error {
	if !n.IsMaster && n.AdvertiseAddr == "" {
		return errors.New("invalid AdvertiseAddr")
	}

	mqttMasterClient := clusterpb.NewMqttMasterClient(n.AdvertiseAddr)

	retryTimes := n.RetryTimes

	for retryTimes > 0 {
		err := mqttMasterClient.Connect()
		if err == nil {
			break
		}

		time.Sleep(n.RetryInterval)
		logx.Info("try connect again, retryTimes :", retryTimes)

		retryTimes--
	}

	_, err := mqttMasterClient.MonitorHandler(context.Background(), &proto.MonitorHandlerRequest{
		CallBackHandler: func(action proto.MonitorAction, serverInfos []*proto.ClusterServerInfo) { // 收到master推送的消息变更

			switch action {
			case proto.MonitorAction_addServer:
				for i := 0; i < len(serverInfos); i++ {
					n.handler.addRemoteService(serverInfos[i])
				}

			case proto.MonitorAction_removeServer:
				for i := 0; i < len(serverInfos); i++ {
					n.handler.delMember("")
				}
			case proto.MonitorAction_replaceServer:

			case proto.MonitorAction_startOve:
			}

		},
	})
	if err != nil {
		return err
	}

	n.masterClient = mqttMasterClient
	return nil
}

// Enable current server accept connection
func (n *Node) listenAndServe() {

}
