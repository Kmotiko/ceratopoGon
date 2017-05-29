package ceratopogon

import (
	"github.com/KMotiko/ceratopogon/messages"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"net"
	"strconv"
	"sync"
)

type MqttSnSession struct {
	mutex    sync.RWMutex
	ClientId string
	Conn     *net.UDPConn
	Remote   *net.UDPAddr
	Topics   *TopicMap
	msgId    *ManagedId

	// duration uint16
	// state bool
}

type TransportSnSession struct {
	// mutex    sync.RWMutex
	MqttSnSession
	mqttClient MQTT.Client
	// brokerAddr string
	// brokerPort int
}

func NewMqttSnSession(id string, conn *net.UDPConn, remote *net.UDPAddr) *MqttSnSession {
	s := &MqttSnSession{
		sync.RWMutex{}, id, conn, remote, NewTopicMap(), &ManagedId{}}
	return s
}

func (s *MqttSnSession) StoreTopic(topicName string) uint16 {
	return s.Topics.StoreTopic(topicName)
}

func (s *MqttSnSession) StoreTopicWithId(topicName string, id uint16) bool {
	return s.Topics.StoreTopicWithId(topicName, id)
}

func (s *MqttSnSession) LoadTopic(topicId uint16) (string, bool) {
	return s.Topics.LoadTopic(topicId)
}

func (s *MqttSnSession) LoadTopicId(topicName string) (uint16, bool) {
	return s.Topics.LoadTopicId(topicName)
}

func (s *MqttSnSession) FreeMsgId(i uint16) {
	s.msgId.FreeId(i)
}

func (s *MqttSnSession) NextMsgId() uint16 {
	return s.msgId.NextId()
}

func NewTransportSnSession(id string, conn *net.UDPConn, remote *net.UDPAddr) *TransportSnSession {
	s := &TransportSnSession{
		MqttSnSession{
			sync.RWMutex{},
			id,
			conn,
			remote,
			NewTopicMap(),
			&ManagedId{}},
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
	log.Println("on publish. Receive message from broker.")
	// get subscribers
	topic := msg.Topic()

	// TODO: check TransportSnSession's state.
	// if session is sleep, gateway must buffer the message.

	// get topicid
	topicId, ok := s.LoadTopicId(topic)

	// if not found
	var m *message.Publish
	msgId := s.NextMsgId()
	if !ok {
		// TODO: implement

		// wildcarded or short topic name

	} else {
		// process as fixed topicid
		// qos, retain, topicId, msgId, data
		// Now, qos is hard coded as 1
		m = message.NewPublishNormal(
			1, false, topicId, msgId, msg.Payload())
	}

	// send message
	s.Conn.WriteToUDP(m.Marshall(), s.Remote)
}
