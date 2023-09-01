package cluster

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"pomelo-go/cluster/clusterpb"
	"pomelo-go/cluster/clusterpb/proto"
	"pomelo-go/component"
	"pomelo-go/tool"
	"time"
)

// Options contains some configurations for current node
type Options struct {
	IsMaster   bool
	ServerId   string                  // node服务id名称
	ServerInfo proto.ClusterServerInfo // node 服务信息用于向master注册

	AdvertiseAddr string        // node服务对应的master地址
	RetryInterval time.Duration // master 重试间隔 default 3*time.sec
	RetryTimes    int           // master 重试间隔 default 10次

	Token string // master 通信token

	Components *component.Components
	//RemoteServiceRoute CustomerRemoteServiceRoute
}

type Node struct {
	Options            // current node options
	ServiceAddr string // current server service address (RPC)

	handler      *LocalHandler          // 处理本地或远程Handler调用
	masterClient clusterpb.MasterClient // 与master通信的客户端 对应pomelo的monitor
	rpcClient    *rpcClient
	server       *clusterpb.MqttMemberServer // 本地rpc服务端

	//sessions map[int64]*session.Session
}

func (n *Node) Startup() error {
	if n.ServiceAddr == "" {
		return errors.New("service address cannot be empty in master node")
	}

	n.rpcClient = newRPCClient()
	n.handler = NewHandler(n)

	components := n.Components.List()
	for _, c := range components {
		err := n.handler.register(c.Comp, c.Opts)
		if err != nil {
			return err
		}
	}

	if err := n.initNode(); err != nil {
		return err
	}

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

// RemoteProcess 远程调用
func (n *Node) RemoteProcess(ctx context.Context, in *proto.RequestRequest) (*proto.RequestResponse, error) {
	return n.handler.remoteProcess(ctx, in)
}

func (n *Node) Shutdown() {
	// reverse call `BeforeShutdown` hooks
	components := n.Components.List()
	length := len(components)
	for i := length - 1; i >= 0; i-- {
		components[i].Comp.BeforeShutdown()
	}
	// reverse call `Shutdown` hooks
	for i := length - 1; i >= 0; i-- {
		components[i].Comp.Shutdown()
	}

	//_, err = client.Unregister(context.Background(), request)

}

func (n *Node) RequestHandler(ctx context.Context, in *proto.RequestRequest) (*proto.RequestResponse, error) {
	logx.Info("node RequestHandler,in:", tool.SimpleJson(in))

	res := []interface{}{
		map[string]interface{}{
			"A": "a",
			"B": "b",
		},
	}

	r := proto.RequestResponse(res)
	return &r, nil
}

func (n *Node) NotifyHandler(ctx context.Context, in *proto.NotifyRequest) {
	logx.Info("node NotifyHandler,in:", tool.SimpleJson(in))

}

func (n *Node) initNode() error {
	if !n.IsMaster && n.AdvertiseAddr == "" {
		return errors.New("invalid AdvertiseAddr")
	}

	n.server = clusterpb.NewMqttMasterServer(n)

	err := n.server.Listen(n.ServiceAddr)
	if err != nil {
		return err
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

	_, err = mqttMasterClient.MonitorHandler(context.Background(), &proto.MonitorHandlerRequest{
		CallBackHandler: func(action proto.MonitorAction, serverInfos []proto.ClusterServerInfo) { // 收到master推送的消息变更

			switch action {
			case proto.MonitorAction_addServer:
				for i := 0; i < len(serverInfos); i++ {

					remoteService, err := transformRemoteServiceInfo(serverInfos[i])
					if err != nil {
						logx.Error("transformRemoteServiceInfo failed,err:", err)
						continue
					}
					n.handler.addRemoteService(remoteService)
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

	_, err = mqttMasterClient.Register(context.Background(), &proto.RegisterRequest{
		ServerInfo: n.ServerInfo,
		Token:      n.Token,
	})
	if err != nil {
		return err
	}

	// 获取注册信息
	subscribeResponse, err := mqttMasterClient.Subscribe(context.Background(), &proto.SubscribeRequest{
		Id: n.ServerId,
	})

	// 初始化handler
	rs := make([]RemoteServiceInfo, 0, len(*subscribeResponse))
	for _, info := range *subscribeResponse {

		remoteService, err := transformRemoteServiceInfo(info)
		if err != nil {
			logx.Error("transformRemoteServiceInfo failed,err:", err)
			continue
		}

		rs = append(rs, remoteService)
	}

	n.handler.initRemoteService(rs)

	n.masterClient = mqttMasterClient
	return nil
}

// Enable current server accept connection
// 与pomelo端通信
func (n *Node) listenAndServe() {

}

func transformRemoteServiceInfo(info proto.ClusterServerInfo) (res RemoteServiceInfo, err error) {

	var (
		host       string
		port       int
		serverType string
	)

	if v, ok := info["host"]; !ok {
		return RemoteServiceInfo{}, errors.New("invalid host")
	} else {
		host = v.(string)
	}

	if v, ok := info["port"]; !ok {
		return RemoteServiceInfo{}, errors.New("invalid port")
	} else {
		port = int(v.(float64))
	}
	if v, ok := info["serverType"]; !ok {
		return RemoteServiceInfo{}, errors.New("invalid serverType")
	} else {
		serverType = v.(string)
	}

	return RemoteServiceInfo{
		ClusterServerInfo: info,
		ServerType:        serverType,
		ServiceAddr:       fmt.Sprintf("%s:%d", host, port),
	}, nil

}
