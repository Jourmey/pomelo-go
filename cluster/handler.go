package cluster

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
	"pomelo-go/cluster/clusterpb/proto"
	"pomelo-go/component"
)

//type rpcHandler func(session *session.Session, msg *message.Message, noCopy bool)

//// CustomerRemoteServiceRoute customer remote service route
//type CustomerRemoteServiceRoute func(service string, session *session.Session, members []*clusterpb.MemberInfo) *clusterpb.MemberInfo

type LocalHandler struct {
	localServices map[string]*component.Service // all registered service
	localHandlers map[string]*component.Handler // all handler method

	mu             sync.RWMutex
	remoteServices map[string][]RemoteServiceInfo // key:serverType value: node serverInfo 数组
	currentNode    *Node
}

type RemoteServiceInfo struct {
	ClusterServerInfo proto.ClusterServerInfo
	ServerType        string
	ServiceAddr       string
}

func NewHandler(currentNode *Node) *LocalHandler {
	h := &LocalHandler{
		localServices:  make(map[string]*component.Service),
		localHandlers:  make(map[string]*component.Handler),
		mu:             sync.RWMutex{},
		remoteServices: make(map[string][]RemoteServiceInfo),
		currentNode:    currentNode,
	}

	return h
}

func (h *LocalHandler) register(comp component.Component, opts []component.Option) error {
	s := component.NewService(comp, opts)

	if _, ok := h.localServices[s.Name]; ok {
		return fmt.Errorf("handler: service already defined: %s", s.Name)
	}

	// register all localHandlers
	h.localServices[s.Name] = s
	return nil
}

func (h *LocalHandler) initRemoteService(members []RemoteServiceInfo) {
	for _, m := range members {
		h.addRemoteService(m)
	}
}

func (h *LocalHandler) addRemoteService(serverInfo RemoteServiceInfo) {
	h.mu.Lock()
	defer h.mu.Unlock()

	logx.Infof("Register remote serverType:%s, serviceAddr:%s", serverInfo.ServerType, serverInfo.ServiceAddr)
	h.remoteServices[serverInfo.ServerType] = append(h.remoteServices[serverInfo.ServerType], serverInfo)
}

func (h *LocalHandler) delMember(addr string) {
	// TODO
}

func (h *LocalHandler) remoteProcess(ctx context.Context, in *proto.RequestRequest) (*proto.RequestResponse, error) {

	// 		Namespace:  "user",
	//		ServerType: "chat",
	//		Service:    "chatRemote",
	//		Method:     "add",

	members := h.findMembers(in.ServerType)
	if len(members) == 0 {

		route := fmt.Sprintf("%s.%s.%s.%s", in.Namespace, in.ServerType, in.Service, in.Method)
		return nil, errors.New(fmt.Sprintf("nano/handler: %s not found(forgot registered?)", route))
	}

	// Select a remote service address
	// 1. if exist customer remote service route ,use it, otherwise use default strategy
	// 2. Use the service address directly if the router contains binding item
	// 3. Select a remote service address randomly and bind to router
	var remoteAddr string
	if false { //h.currentNode.Options.RemoteServiceRoute != nil {
		//if addr, found := session.Router().Find(service); found {
		//	remoteAddr = addr
		//} else {
		//	member := h.currentNode.Options.RemoteServiceRoute(service, session, members)
		//	if member == nil {
		//		log.Println(fmt.Sprintf("customize remoteServiceRoute handler: %s is not found", msg.Route))
		//		return
		//	}
		//	remoteAddr = member.ServiceAddr
		//	session.Router().Bind(service, remoteAddr)
		//}
	} else {

		remoteAddr = members[rand.Intn(len(members))].ServiceAddr

	}

	pool, err := h.currentNode.rpcClient.getConnPool(remoteAddr)
	if err != nil {
		return nil, err
	}

	client := pool.Get()

	return client.Request(ctx, in)
}

func (h *LocalHandler) findMembers(service string) []RemoteServiceInfo {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.remoteServices[service]
}

func (h *LocalHandler) localProcess() {

}
