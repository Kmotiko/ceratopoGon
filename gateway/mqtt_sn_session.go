package ceratopogon

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"net"
	"strconv"
)

type MqttSnSession struct {
	ClientId string
	Conn     *net.UDPConn
	Remote   *net.UDPAddr
	TopicMap map[uint16]string
	topicId  uint16
	msgId    [0xffff]bool
}

type TransportSnSession struct {
	MqttSnSession
	mqttClient MQTT.Client
	// brokerAddr string
	// brokerPort int
}

func NewMqttSnSession(id string, conn *net.UDPConn, remote *net.UDPAddr) *MqttSnSession {
	s := &MqttSnSession{
		id, conn, remote, make(map[uint16]string), 0, [0xffff]bool{}}
	return s
}

func (s *MqttSnSession) NextTopicId() uint16 {
	// TODO: add lock
	s.topicId++
	if s.topicId == 0xffff {
		// err
	}

	return s.topicId
}

func (s *MqttSnSession) FreeMsgId(i uint16) {
	// TODO: add lock
	s.msgId[i] = false
}

func (s *MqttSnSession) NextMsgId() uint16 {
	// TODO: add lock
	for i, v := range s.msgId {
		if !v {
			return uint16(i)
		}
	}
	// TODO: error handling

	return 0
}

func NewTransportSnSession(id string, conn *net.UDPConn, remote *net.UDPAddr) *TransportSnSession {
	s := &TransportSnSession{
		MqttSnSession{
			id,
			conn,
			remote,
			make(map[uint16]string),
			0,
			[0xffff]bool{}},
		nil}
	return s
}

func (s *TransportSnSession) ConnectToBroker(brokerAddr string, brokerPort int, user string, password string) error {
	// create opts
	addr := "tcp://" + brokerAddr + ":" + strconv.Itoa(brokerPort)
	opts := MQTT.NewClientOptions().AddBroker(addr)
	opts.SetClientID(s.MqttSnSession.ClientId)
	opts.SetUsername(user)
	opts.SetPassword(password)

	// create client instance
	s.mqttClient = (MQTT.NewClient(opts))

	// connect
	if token := s.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

func (s *TransportSnSession) OnPublish(client MQTT.Client, msg MQTT.Message) {
	// send Publish to MQTT-SN Client
	log.Println("on publish!!!")
}
