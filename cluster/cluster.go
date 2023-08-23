package cluster

import (
	"sync"
)

type cluster struct {
	// If cluster is not large enough, use slice is OK
	currentNode *Node
	rpcClient   *rpcClient

	mu      sync.RWMutex
	members []*Member
}

func newCluster(currentNode *Node) *cluster {
	c := &cluster{currentNode: currentNode}
	if currentNode.IsMaster {
		c.checkMemberHeartbeat()
	}
	return c
}

func (c *cluster) setRpcClient(client *rpcClient) {
	c.rpcClient = client
}

func (c *cluster) remoteAddrs() []string {
	var addrs []string
	c.mu.RLock()
	for _, m := range c.members {
		addrs = append(addrs, m.memberInfo.ServiceAddr)
	}
	c.mu.RUnlock()
	return addrs
}

func (c *cluster) initMembers(members []*clusterpb.MemberInfo) {
	c.mu.Lock()
	for _, info := range members {
		c.members = append(c.members, &Member{
			memberInfo: info,
		})
	}
	c.mu.Unlock()
}

func (c *cluster) addMember(info *clusterpb.MemberInfo) {
	c.mu.Lock()
	var found bool
	for _, member := range c.members {
		if member.memberInfo.ServiceAddr == info.ServiceAddr {
			member.memberInfo = info
			found = true
			break
		}
	}
	if !found {
		c.members = append(c.members, &Member{
			memberInfo: info,
		})
	}
	c.mu.Unlock()
}

func (c *cluster) delMember(addr string) {
	c.mu.Lock()
	var index = -1
	for i, member := range c.members {
		if member.memberInfo.ServiceAddr == addr {
			index = i
			break
		}
	}
	if index != -1 {
		if index >= len(c.members)-1 {
			c.members = c.members[:index]
		} else {
			c.members = append(c.members[:index], c.members[index+1:]...)
		}
	}
	c.mu.Unlock()
}
