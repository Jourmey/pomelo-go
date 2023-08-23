package clusterpb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"strconv"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"pomelo-go/cluster/clusterpb/proto"
)

const (
	topic_Register = "register"
	topic_Monitor  = "monitor"

	action_Subscribe = "subscribe"

	pro_ok   = 1
	pro_fail = -1
)

type requestCallback func(msg monitorMessage)

type MqttMasterClient struct {
	clientId string // = 'MQTT_ADMIN_' + Date.now();

	host           string
	port           int
	keepaliveTimer time.Duration // default 2s
	pingTimeout    time.Duration // default 1s

	reqId  int
	socket mqtt.Client
	resp   sync.Map // 返回值队列
	//callbacks sync.Map //  map[int]requestCallback

	register  chan registerResponse
	subscribe chan proto.ClusterServerInfo
}

func (m *MqttMasterClient) Register(ctx context.Context, in *proto.RegisterRequest) (*proto.RegisterResponse, error) {

	req := make(map[string]interface{}, len(in.ServerInfo)+1)

	for s, i := range in.ServerInfo {
		req[s] = i
	}
	req["token"] = in.Token

	err := m.doSend(topic_Register, req)
	if err != nil {
		return nil, err
	}

	select {
	case res := <-m.register:

		if res.Code == pro_ok {
			return &proto.RegisterResponse{}, nil
		}

		return nil, errors.New(res.Msg)

	case <-time.After(10 * time.Second):
		return nil, errors.New("receive register timeout")
	}

}

func (m *MqttMasterClient) Subscribe(ctx context.Context, in *proto.SubscribeRequest) (*proto.SubscribeResponse, error) {

	request := subscribeRequest{
		Action: action_Subscribe,
		Id:     in.Id,
	}

	serverInfos := make(map[string]proto.ClusterServerInfo)

	err := m.Request(proto.MASTER_WATCHER, request, func(msg monitorMessage) {

		for serverId, serverInfo := range msg.Body {

			si, ok := serverInfo.(map[string]interface{})
			if !ok {
				logx.Errorf("Subscribe serverInfo.(map[string]interface{}) failed ")
				continue
			}

			serverInfos[serverId] = proto.ClusterServerInfo(si)
		}

		logx.Info("Subscribe success")
	})

	res := proto.SubscribeResponse(serverInfos)
	return &res, err
}

func (m *MqttMasterClient) Record(ctx context.Context, in *proto.RecordRequest) (*proto.RecordResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MqttMasterClient) Connect() error {

	token := m.socket.Connect()

	token.Wait()

	return token.Error()
}

func (m *MqttMasterClient) publishHandler(client mqtt.Client, message mqtt.Message) {

	switch message.Topic() {

	case topic_Register:

		res := registerResponse{}
		err := json.Unmarshal(message.Payload(), &res)
		if err != nil {
			return
		}

		select {
		case m.register <- res:
		default:
			logx.Errorf("topic_Register chan failed")
		}

	case topic_Monitor:

		msg := monitorMessage{}

		// 这里接收的字符串居然是转义后的
		unescapedString, err := strconv.Unquote(string(message.Payload()))
		if err != nil {
			return
		}

		err = json.Unmarshal([]byte(unescapedString), &msg)
		if err != nil {
			return
		}

		if msg.Command != nil {

		} else if msg.RespId != nil {

		} else {

		}

	default:

		logx.Errorf("invalid topic")

	}

}

func (m *MqttMasterClient) Notify(moduleId string, body interface{}) error {
	return m.doSend(topic_Monitor, map[string]interface{}{

		"moduleId": moduleId,
		"body":     body,
	})
}

func (m *MqttMasterClient) Request(moduleId string, body interface{}, cb requestCallback) error {

	m.reqId++
	var reqId = m.reqId
	err := m.doSend(topic_Monitor, map[string]interface{}{
		"reqId":    reqId,
		"moduleId": moduleId,
		"body":     body,
	})

	if err != nil {
		return err
	}

	m.callbacks.Store(reqId, cb)

	return nil
}

func (m *MqttMasterClient) doSend(topic string, msg interface{}) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if pToken := m.socket.Publish(topic, 0, false, payload); pToken.Wait() && pToken.Error() != nil {
		return pToken.Error()
	}

	return nil
}

func NewMasterClient(host string, port int) *MqttMasterClient {

	var (
		clientId       = fmt.Sprintf("MQTT_ADMIN_%d", time.Now().UnixMilli())
		keepaliveTimer = 2 * time.Second
		pingTimeout    = 1 * time.Second
	)

	m := &MqttMasterClient{
		clientId:       clientId,
		host:           host,
		port:           port,
		keepaliveTimer: keepaliveTimer,
		pingTimeout:    pingTimeout,
		socket:         nil,
		register:       make(chan registerResponse),
		subscribe:      make(chan proto.ClusterServerInfo),
	}

	opts := mqtt.NewClientOptions().
		AddBroker(fmt.Sprintf("tcp://%s:%d", m.host, m.port)).
		SetClientID(m.clientId)

	opts.SetKeepAlive(m.keepaliveTimer)
	opts.SetDefaultPublishHandler(m.publishHandler)
	opts.SetPingTimeout(m.pingTimeout)

	socket := mqtt.NewClient(opts)
	m.socket = socket

	return m
}

type monitorMessage struct {
	RespId *int    `json:"respId"` //  "respId": 1,
	Error  *string `json:"error"`  //  "error": null,

	ReqId    *int    `json:"reqId"`    //  "reqId": 1,
	ModuleId *string `json:"moduleId"` //  "moduleId": "__monitorwatcher__",

	Command *string `json:"command"` // command

	Body map[string]interface{} `json:"body"` //  "body": {
}

type registerResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type subscribeRequest struct {
	Action string `json:"action"`
	Id     string `json:"id"`
}
