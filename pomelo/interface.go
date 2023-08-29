package pomelo

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
	"os/signal"
	"pomelo-go/cluster"
	"pomelo-go/cluster/clusterpb/proto"
	"pomelo-go/component"
	"sync/atomic"
	"syscall"
	"time"
)

var (
	running     int32
	currentNode *cluster.Node
	die         chan bool // wait for end application
)

func init() {
	die = make(chan bool)
}

func Listen(addr string, opts ...Option) {
	if atomic.AddInt32(&running, 1) != 1 {
		logx.Info("Nano has running")
		return
	}

	opt := cluster.Options{
		Components: &component.Components{},
	}
	for _, option := range opts {
		option(&opt)
	}

	// Set the retry interval to 3 secondes if doesn't set by user
	if opt.RetryInterval == 0 {
		opt.RetryInterval = time.Second * 3
	}

	if opt.RetryTimes == 0 {
		opt.RetryTimes = 10
	}

	node := &cluster.Node{
		Options:     opt,
		ServiceAddr: addr,
	}

	err := node.Startup()
	if err != nil {
		logx.Infof("Node startup failed: %v", err)
	}
	currentNode = node

	logx.Infof(fmt.Sprintf("Startup *Nano backend server* %s, service address %s",
		node.ServerId, node.ServiceAddr))

	sg := make(chan os.Signal)
	signal.Notify(sg, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM)

	select {
	case <-die:
		logx.Info("The app will shutdown in a few seconds")
	case s := <-sg:
		logx.Info("Nano server got signal", s)
	}

	logx.Info("Nano server is stopping...")

	node.Shutdown()
	currentNode = nil

	atomic.StoreInt32(&running, 0)
}

// Shutdown send a signal to let 'nano' shutdown itself.
func Shutdown() {
	close(die)
}

func RemoteProcess(ctx context.Context, in *proto.RequestRequest) (*proto.RequestResponse, error) {
	if currentNode == nil {
		return nil, errors.New("invalid node")
	}

	return currentNode.RemoteProcess(ctx, in)
}
