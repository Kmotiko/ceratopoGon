package ceratopoGon

import (
	"github.com/Kmotiko/ceratopoGon/messages"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"net"
	"strconv"
	"sync"
)

type MqttSnSession struct {
	mutex        sync.RWMutex
	ClientId     string
	Conn         *net.UDPConn
	Remote       *net.UDPAddr
	Topics       *TopicMap
	PredefTopics *TopicMap
	msgId        *ManagedId

	// duration uint16
	// state bool
}

type TransportSnSession struct {
	// mutex    sync.RWMutex
	*MqttSnSession
	mqttClient MQTT.Client
	// brokerAddr string
	// brokerPort int
}

func NewMqttSnSession(id string, conn *net.UDPConn, remote *net.UDPAddr) *MqttSnSession {
	s := &MqttSnSession{
		sync.RWMutex{},
		id,
		conn,
		remote,
		NewTopicMap(),
		NewTopicMap(),
		&ManagedId{}}
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

func (s *MqttSnSession) StorePredefTopic(topicName string) uint16 {
	return s.PredefTopics.StoreTopic(topicName)
}

func (s *MqttSnSession) StorePredefTopicWithId(topicName string, id uint16) bool {
	return s.PredefTopics.StoreTopicWithId(topicName, id)
}

func (s *MqttSnSession) LoadPredefTopic(topicId uint16) (string, bool) {
	return s.PredefTopics.LoadTopic(topicId)
}

func (s *MqttSnSession) LoadPredefTopicId(topicName string) (uint16, bool) {
	return s.PredefTopics.LoadTopicId(topicName)
}
func (s *MqttSnSession) FreeMsgId(i uint16) {
	s.msgId.FreeId(i)
}

func (s *MqttSnSession) NextMsgId() uint16 {
	return s.msgId.NextId()
}

func NewTransportSnSession(id string, conn *net.UDPConn, remote *net.UDPAddr) *TransportSnSession {
	s := &TransportSnSession{
		&MqttSnSession{
			sync.RWMutex{},
			id,
			conn,
			remote,
			NewTopicMap(),
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
	// get topic
	topic := msg.Topic()

	// TODO: check TransportSnSession's state.
	// if session is sleep, gateway must buffer the message.

	var m *message.Publish
	msgId := s.NextMsgId()

	// get predef topicid
	topicId, ok := s.LoadPredefTopicId(topic)
	if ok {
		// process as predef topicid
		// qos, retain, topicId, msgId, data
		// Now, qos is hard coded as 1
		m = message.NewPublishPredefined(
			1, false, topicId, msgId, msg.Payload())
	} else if topicId, ok = s.LoadTopicId(topic); ok {
		// process as fixed topicid
		// qos, retain, topicId, msgId, data
		// Now, qos is hard coded as 1
		m = message.NewPublishNormal(
			1, false, topicId, msgId, msg.Payload())
	} else {
		// TODO: implement wildcarded or short topic route

		log.Println("[Warn] topic id was not found for ", topic, ".")
		return
	}

	// send message
	s.Conn.WriteToUDP(m.Marshall(), s.Remote)
}
