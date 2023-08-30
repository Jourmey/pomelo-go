package clusterpb

import (
	"bytes"
	"context"
	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/zeromicro/go-zero/core/logx"
	"pomelo-go/cluster/clusterpb/proto"
)

type MqttMasterServer struct {
	server *mqtt.Server
}

func (m *MqttMasterServer) RequestHandler(ctx context.Context, in *proto.RequestRequest) (*proto.RequestResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MqttMasterServer) NotifyHandler(ctx context.Context, in *proto.NotifyRequest) (*proto.NotifyResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m *MqttMasterServer) Listener(advertiseAddr string) error {
	server := mqtt.New(nil)
	m.server = server

	ws := listeners.NewTCP("t1", advertiseAddr, nil)
	_ = ws

	err := server.AddHook(new(ExampleHook), nil)
	if err != nil {
		return err
	}

	err = server.AddListener(ws)
	if err != nil {
		return err
	}

	err = server.Serve()
	if err != nil {
		return err
	}

	return nil
}

func NewMqttMasterServer() *MqttMasterServer {

	return &MqttMasterServer{
		server: nil,
	}
}

type ExampleHook struct {
	mqtt.HookBase
}

func (h *ExampleHook) ID() string {
	return "events-example"
}

func (h *ExampleHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnect,
		mqtt.OnConnectAuthenticate,
		mqtt.OnDisconnect,

		mqtt.OnSubscribe,
		mqtt.OnSubscribed,
		mqtt.OnUnsubscribed,
		mqtt.OnSelectSubscribers,

		//mqtt.OnPacketRead,
		//mqtt.OnPacketEncode,
		//mqtt.OnPacketSent,
		mqtt.OnPacketProcessed,

		mqtt.OnPacketProcessed,

		mqtt.OnPublished,
		mqtt.OnPublish,
	}, []byte{b})
}

func (h *ExampleHook) Init(config any) error {
	h.Log.Info().Msg("initialised")
	return nil
}

func (h *ExampleHook) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	// 全部放行
	logx.Infof("OnConnectAuthenticate , id:%s", cl.ID)
	return true
}

func (h *ExampleHook) OnConnect(cl *mqtt.Client, pk packets.Packet) error {
	h.Log.Info().Str("client", cl.ID).Msgf("client connected")
	return nil
}

func (h *ExampleHook) OnDisconnect(cl *mqtt.Client, err error, expire bool) {
	h.Log.Info().Str("client", cl.ID).Bool("expire", expire).Err(err).Msg("client disconnected")
}

func (h *ExampleHook) OnSubscribe(cl *mqtt.Client, pk packets.Packet) packets.Packet {
	logx.Infof("OnSubscribe , id:%s", cl.ID)
	return pk
}

func (h *ExampleHook) OnSubscribed(cl *mqtt.Client, pk packets.Packet, reasonCodes []byte) {
	h.Log.Info().Str("client", cl.ID).Interface("filters", pk.Filters).Msgf("subscribed qos=%v", reasonCodes)

}

func (h *ExampleHook) OnUnsubscribed(cl *mqtt.Client, pk packets.Packet) {
	h.Log.Info().Str("client", cl.ID).Interface("filters", pk.Filters).Msg("unsubscribed")
}

func (h *ExampleHook) OnSelectSubscribers(subs *mqtt.Subscribers, pk packets.Packet) *mqtt.Subscribers {
	logx.Infof("OnSelectSubscribers")
	return subs
}

// triggers after a packet from the client been processed (handled)
func (h *ExampleHook) OnPacketProcessed(cl *mqtt.Client, pk packets.Packet, err error) {
	if err != nil {
		logx.Errorf("OnPacketProcessed failed , id:%s , err:%s", cl.ID, err)
	} else {
		logx.Infof("OnPacketProcessed , id:%s, pk:%+v,err:%s", cl.ID, pk, err)
	}

	switch pk.FixedHeader.Type {

	case packets.Reserved:
		logx.Infof("OnPacketProcessed Reserved , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Connect:
		logx.Infof("OnPacketProcessed Connect , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Connack:
		logx.Infof("OnPacketProcessed Connack , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Publish:
		logx.Infof("OnPacketProcessed Publish , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Puback:
		logx.Infof("OnPacketProcessed Puback , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Pubrec:
		logx.Infof("OnPacketProcessed Pubrec , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Pubrel:
		logx.Infof("OnPacketProcessed Pubrel , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Pubcomp:
		logx.Infof("OnPacketProcessed Pubcomp , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Subscribe:
		logx.Infof("OnPacketProcessed Subscribe , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Suback:
		logx.Infof("OnPacketProcessed Suback , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Unsubscribe:
		logx.Infof("OnPacketProcessed Unsubscribe , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Unsuback:
		logx.Infof("OnPacketProcessed Unsuback , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Pingreq:
		logx.Infof("OnPacketProcessed Pingreq , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Pingresp:
		logx.Infof("OnPacketProcessed Pingresp , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Disconnect:
		logx.Infof("OnPacketProcessed Disconnect , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.Auth:
		logx.Infof("OnPacketProcessed Auth , id:%s, payload:%s", cl.ID, pk.Payload)

	case packets.WillProperties:
		logx.Infof("OnPacketProcessed WillProperties , id:%s, payload:%s", cl.ID, pk.Payload)

	default:
		logx.Infof("OnPacketProcessed default , id:%s, payload:%s", cl.ID, pk.Payload)

	}

}

func (h *ExampleHook) OnPublish(cl *mqtt.Client, pk packets.Packet) (packets.Packet, error) {
	h.Log.Info().Str("client", cl.ID).Str("payload", string(pk.Payload)).Msg("received from client")

	pkx := pk
	if string(pk.Payload) == "hello" {
		pkx.Payload = []byte("hello world")
		h.Log.Info().Str("client", cl.ID).Str("payload", string(pkx.Payload)).Msg("received modified packet from client")
	}

	return pkx, nil
}

func (h *ExampleHook) OnPublished(cl *mqtt.Client, pk packets.Packet) {
	h.Log.Info().Str("client", cl.ID).Str("payload", string(pk.Payload)).Msg("published to client")
}
