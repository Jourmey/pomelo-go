package cluster

import (
	"github.com/zeromicro/go-zero/core/logx"
	"pomelo-go/cluster/clusterpb/proto"
	"pomelo-go/component"
	"sync"
)

//type rpcHandler func(session *session.Session, msg *message.Message, noCopy bool)
//
//// CustomerRemoteServiceRoute customer remote service route
//type CustomerRemoteServiceRoute func(service string, session *session.Session, members []*clusterpb.MemberInfo) *clusterpb.MemberInfo

type LocalHandler struct {
	localServices map[string]*component.Service // all registered service
	localHandlers map[string]*component.Handler // all handler method

	mu             sync.RWMutex
	remoteServices map[string][]*proto.ClusterServerInfo // key:serverType value: node serverInfo 数组
	currentNode    *Node
}

func NewHandler(currentNode *Node) *LocalHandler {
	h := &LocalHandler{
		localServices:  make(map[string]*component.Service),
		localHandlers:  make(map[string]*component.Handler),
		mu:             sync.RWMutex{},
		remoteServices: make(map[string][]*proto.ClusterServerInfo),
		currentNode:    currentNode,
	}

	return h
}

func (h *LocalHandler) register(comp component.Component, opts []component.Option) error {

	return nil
}

func (h *LocalHandler) initRemoteService(members map[string]*proto.ClusterServerInfo) {
	for _, m := range members {
		h.addRemoteService(m)
	}
}

func (h *LocalHandler) addRemoteService(serverInfo *proto.ClusterServerInfo) {
	h.mu.Lock()
	defer h.mu.Unlock()

	logx.Info("Register remote service", serverInfo.ServerType)
	h.remoteServices[serverInfo.ServerType] = append(h.remoteServices[serverInfo.ServerType], serverInfo)
}

func (h *LocalHandler) delMember(addr string) {
	// TODO
}

func (h *LocalHandler) remoteProcess() {

}

func (h *LocalHandler) processMessage() {

	//handler, found := h.localHandlers[msg.Route]
	//if !found {
	//	h.remoteProcess(agent.session, msg, false)
	//} else {
	//	h.localProcess(handler, lastMid, agent.session, msg)
	//}
}

func (h *LocalHandler) localProcess() {

}
