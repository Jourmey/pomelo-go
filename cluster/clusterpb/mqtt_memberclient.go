package clusterpb

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/zeromicro/go-zero/core/logx"
	"pomelo-go/cluster/clusterpb/proto"
)

const (
	topic_RPC = "RPC"
)

type MqttMemberClient struct {
	clientId string // = 'MQTT_RPC_' + Date.now();

	advertiseAddr  string
	keepaliveTimer time.Duration // default 2s
	pingTimeout    time.Duration // default 1s
	requestTimeout time.Duration // default 10s

	reqId  int
	socket mqtt.Client
	resp   sync.Map // monitor memberRequest 请求列表
}

func (m *MqttMemberClient) Request(ctx context.Context, in *proto.RequestRequest) (*proto.RequestResponse, error) {

	type rpcData struct {
		ID  int         `json:"id"`
		Msg interface{} `json:"msg"`
	}

	type message struct {
		Namespace  string        `json:"namespace"`
		ServerType string        `json:"serverType"`
		Service    string        `json:"service"`
		Method     string        `json:"method"`
		Args       []interface{} `json:"args"`
	}

	msg := message{
		Namespace:  "user",
		ServerType: "chat",
		Service:    "chatRemote",
		Method:     "add",
		Args: []interface{}{
			"stu1*kick_testsss",
			"cluster-server-connector-0",
			"kick_testsss",
			true,
			2,
			1,
			0,
			"123",
			"abc",
			"2.9.8.7",
			map[string]interface{}{
				"uniqId":    "231FF2BB-BA09-598D-9EB6-3B0299D292E7ssss",
				"rid":       "kick_testsss",
				"rtype":     2,
				"role":      1,
				"ulevel":    0,
				"uname":     "123",
				"classid":   "abc",
				"clientVer": "2.9.8.7",
				"userVer":   "1.0",
				"liveType":  "COMBINE_SMALL_CLASS_MODE"},
			"0"},
	}

	err := m.doSend(topic_RPC, rpcData{
		ID:  1,
		Msg: msg,
	})

	return &proto.RequestResponse{}, err

}

func (m *MqttMemberClient) Notify(ctx context.Context, in *proto.NotifyRequest) (*proto.NotifyResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MqttMemberClient) Connect() error {

	token := m.socket.Connect()

	token.Wait()

	return token.Error()
}

func (m *MqttMemberClient) doSend(topic string, msg interface{}) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if pToken := m.socket.Publish(topic, 0, false, payload); pToken.Wait() && pToken.Error() != nil {
		return pToken.Error()
	}

	return nil
}

func (m *MqttMemberClient) publishHandler(client mqtt.Client, message mqtt.Message) {

	logx.Debugf("publishHandler,message: %s", message.Payload())

	switch message.Topic() {

	case topic_RPC:

		msg := rpcMessage{}

		//// 这里接收的字符串居然是转义后的
		//unescapedString, err := strconv.Unquote(string(message.Payload()))
		//if err != nil {
		//	return
		//}

		err := json.Unmarshal(message.Payload(), &msg)
		if err != nil {
			logx.Error("MqttMemberClient publishHandler json.Unmarshal failed ,err:", err)
			return
		}

		req, ok := m.resp.LoadAndDelete(*msg.RespId)
		if !ok {
			logx.Error("MqttMemberClient publishHandler LoadAndDelete failed")
			return
		}

		mReq := req.(memberRequest)

		select {
		case mReq.resp <- msg:
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
		SetIgnoreVerifyConnACK(true)

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
	RespId *int    `json:"respId"` //  "respId": 1,
	Error  *string `json:"error"`  //  "error": null,

	ReqId    *int    `json:"reqId"`    //  "reqId": 1,
	ModuleId *string `json:"moduleId"` //  "moduleId": "__monitorwatcher__",

	Command *string `json:"command"` // command

	Body json.RawMessage `json:"body"` // 不同返回值的
}
