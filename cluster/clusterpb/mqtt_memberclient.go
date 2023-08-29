package clusterpb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zeromicro/go-zero/core/logx"
	"pomelo-go/cluster/clusterpb/proto"
)

const (
	topic_RPC = "rpc"
)

type MqttMemberClient struct {
	clientId string // = 'MQTT_RPC_' + Date.now();

	advertiseAddr  string
	keepaliveTimer time.Duration // default 2s
	pingTimeout    time.Duration // default 1s
	requestTimeout time.Duration // default 10s

	reqId      int
	reqIdMutex sync.Mutex

	socket mqtt.Client
	resp   sync.Map // monitor memberRequest 请求列表
}

func (m *MqttMemberClient) Request(ctx context.Context, in *proto.RequestRequest) (*proto.RequestResponse, error) {

	type rpcRequestData struct {
		ID  int         `json:"id"`
		Msg interface{} `json:"msg"`
	}

	m.reqIdMutex.Lock()
	m.reqId++
	var reqId = m.reqId
	m.reqIdMutex.Unlock()

	err := m.doSend(topic_RPC, rpcRequestData{
		ID:  reqId,
		Msg: in,
	})

	if err != nil {
		return nil, err
	}

	r := memberRequest{
		resp:  make(chan rpcMessage),
		reqId: reqId,
	}

	m.resp.Store(reqId, r)

	select {
	case resp := <-r.resp:
		res := proto.RequestResponse{}

		err := json.Unmarshal(resp.Resp, &res)
		if err != nil {
			return nil, err
		}

		return &res, nil

	case <-time.After(m.requestTimeout):
		return nil, errors.New("timeout")
	}

}

func (m *MqttMemberClient) Notify(ctx context.Context, in *proto.NotifyRequest) (*proto.NotifyResponse, error) {
	type rpcNotifyData struct {
		Msg interface{} `json:"msg"`
	}

	err := m.doSend(topic_RPC, rpcNotifyData{
		Msg: in,
	})

	if err != nil {
		return nil, err
	}

	return &proto.NotifyResponse{}, nil
}

func (m *MqttMemberClient) Connect() error {

	token := m.socket.Connect()

	token.Wait()

	return token.Error()
}

func (m *MqttMemberClient) Close() error {
	return nil
}

func (m *MqttMemberClient) doSend(topic string, msg interface{}) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	logx.Debugf("send message: %s", payload)

	if pToken := m.socket.Publish(topic, 0, false, payload); pToken.Wait() && pToken.Error() != nil {
		return pToken.Error()
	}

	return nil
}

func (m *MqttMemberClient) publishHandler(client mqtt.Client, message mqtt.Message) {

	logx.Debugf("receive message: %s", message.Payload())

	switch message.Topic() {

	case topic_RPC:

		msg := rpcMessage{}

		err := json.Unmarshal(message.Payload(), &msg)
		if err != nil {
			logx.Error("MqttMemberClient publishHandler json.Unmarshal failed ,err:", err)
			return
		}

		req, ok := m.resp.LoadAndDelete(msg.Id)
		if !ok {
			logx.Error("MqttMemberClient publishHandler LoadAndDelete failed")
			return
		}

		mReq := req.(memberRequest)

		select {
		case mReq.resp <- msg:
			close(mReq.resp)
		default:
			logx.Error("monitorRequest chan failed")
		}

	default:

		logx.Error("invalid topic")

	}

}

func NewMqttMemberClient(advertiseAddr string) *MqttMemberClient {

	var (
		clientId       = fmt.Sprintf("MQTT_RPC_%d", time.Now().UnixMilli())
		keepaliveTimer = 2 * time.Second
		pingTimeout    = 1 * time.Second
		requestTimeout = 5 * time.Second
	)

	m := &MqttMemberClient{
		clientId:       clientId,
		advertiseAddr:  advertiseAddr,
		keepaliveTimer: keepaliveTimer,
		pingTimeout:    pingTimeout,
		requestTimeout: requestTimeout,
		reqId:          0,
		socket:         nil,
		resp:           sync.Map{},
	}

	opts := mqtt.NewClientOptions().
		AddBroker(advertiseAddr).
		SetClientID(m.clientId).
		SetCleanSession(false).
		SetIgnoreVerifyConnACK(true) // 这里对mqtt做了适配改造，pomelo的服务端不会回复connACK，正常逻辑导致连接失败

	//opts.SetKeepAlive(m.keepaliveTimer)
	opts.SetDefaultPublishHandler(m.publishHandler)
	opts.SetPingTimeout(m.pingTimeout)

	socket := mqtt.NewClient(opts)
	m.socket = socket

	return m
}

type memberRequest struct {
	resp  chan rpcMessage
	reqId int
}

type rpcMessage struct {
	Id   int             `json:"id"`   //  "respId": 1,
	Resp json.RawMessage `json:"resp"` // 不同返回值的
}
