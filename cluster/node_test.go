package cluster

import "testing"

func TestNode_Startup(t *testing.T) {

	n := &Node{
		Options: Options{
			IsMaster:           false,
			AdvertiseAddr:      "",
			Components:         nil,
			RemoteServiceRoute: nil,
		},

		ServiceAddr: "127.0.0.1:4450",
		cluster:     nil,
		handler:     nil,
		server:      nil,
		rpcClient:   nil,
		sessions:    nil,
	}

}
