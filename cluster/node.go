package cluster

import "pomelo-go/component"

// Options contains some configurations for current node
type Options struct {
	IsMaster           bool
	AdvertiseAddr      string
	Components         *component.Components
	RemoteServiceRoute CustomerRemoteServiceRoute
}

type Node struct {
	Options            // current node options
	ServiceAddr string // current server service address (RPC)

	cluster *cluster
	handler *LocalHandler
	monitor *monitorClient // 与master通信的客户端
	server  *grpcServer    // rpc服务端

	sessions map[int64]*session.Session
}

func (n *Node) Startup() error {

	return nil
}

func (n *Node) Handler() *LocalHandler {
	return n.handler
}

func (n *Node) initNode() error {

	return nil
}

func (n *Node) Shutdown() {

}

// Enable current server accept connection
func (n *Node) listenAndServe() {

}
