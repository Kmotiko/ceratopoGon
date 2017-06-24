package ceratopoGon

import (
	"github.com/Kmotiko/ceratopoGon/messages"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

/**
 *
 */
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

/**
 *
 */
type TransportSnSession struct {
	mutex sync.RWMutex
	*MqttSnSession
	mqttClient     MQTT.Client
	brokerHost     string
	brokerPort     int
	brokerUser     string
	brokerPassword string
	sendBuffer     chan message.MqttSnMessage
	shutDown       chan bool
}

/**
 *
 */
func NewMqttSnSession(id string, conn *net.UDPConn, remote *net.UDPAddr) *MqttSnSession {
	s := &MqttSnSession{
		sync.RWMutex{},
		id,
		conn,
		remote,
		NewTopicMap(),
		NewTopicMap(),
		&ManagedId{},
	}
	return s
}

/**
 *
 */
func (s *MqttSnSession) StoreTopic(topicName string) uint16 {
	return s.Topics.StoreTopic(topicName)
}

/**
 *
 */
func (s *MqttSnSession) StoreTopicWithId(topicName string, id uint16) bool {
	return s.Topics.StoreTopicWithId(topicName, id)
}

/**
 *
 */
func (s *MqttSnSession) LoadTopic(topicId uint16) (string, bool) {
	return s.Topics.LoadTopic(topicId)
}

/**
 *
 */
func (s *MqttSnSession) LoadTopicId(topicName string) (uint16, bool) {
	return s.Topics.LoadTopicId(topicName)
}

/**
 *
 */
func (s *MqttSnSession) StorePredefTopic(topicName string) uint16 {
	return s.PredefTopics.StoreTopic(topicName)
}

/**
 *
 */
func (s *MqttSnSession) StorePredefTopicWithId(topicName string, id uint16) bool {
	return s.PredefTopics.StoreTopicWithId(topicName, id)
}

/**
 *
 */
func (s *MqttSnSession) LoadPredefTopic(topicId uint16) (string, bool) {
	return s.PredefTopics.LoadTopic(topicId)
}

/**
 *
 */
func (s *MqttSnSession) LoadPredefTopicId(topicName string) (uint16, bool) {
	return s.PredefTopics.LoadTopicId(topicName)
}

/**
 *
 */
func (s *MqttSnSession) FreeMsgId(i uint16) {
	s.msgId.FreeId(i)
}

/**
 *
 */
func (s *MqttSnSession) NextMsgId() uint16 {
	return s.msgId.NextId()
}

/**
 *
 */
func NewTransportSnSession(
	id string,
	conn *net.UDPConn,
	remote *net.UDPAddr,
	host string,
	port int,
	user string,
	password string) *TransportSnSession {
	s := &TransportSnSession{
		sync.RWMutex{},
		&MqttSnSession{
			sync.RWMutex{},
			id,
			conn,
			remote,
			NewTopicMap(),
			NewTopicMap(),
			&ManagedId{}},
		nil, host, port, user, password,
		make(chan message.MqttSnMessage, 10),
		make(chan bool, 1)}
	return s
}

/**
 *
 */
func (s *TransportSnSession) connLostHandler(
	c MQTT.Client, err error) {
	log.Println("ERROR : MQTT connection is lost with ", err)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	err = s.ConnectToBroker(false)
	if err != nil {
		log.Println("ERROR : failed to connect to broker")
		os.Exit(0)
	}
}

/**
 *
 */
func (s *TransportSnSession) ConnectToBroker(cleanSession bool) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// create opts
	addr := "tcp://" + s.brokerHost + ":" + strconv.Itoa(s.brokerPort)
	opts := MQTT.NewClientOptions().AddBroker(addr)
	opts.SetClientID(s.MqttSnSession.ClientId)
	opts.SetUsername(s.brokerUser)
	opts.SetPassword(s.brokerPassword)
	opts.SetAutoReconnect(false)
	opts.SetConnectionLostHandler(
		s.connLostHandler)
	opts.SetCleanSession(cleanSession)

	// create client instance
	s.mqttClient = (MQTT.NewClient(opts))

	// connect
	if token := s.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

/**
 *
 */
func (s *TransportSnSession) OnPublish(client MQTT.Client, msg MQTT.Message) {
	log.Println("on publish. Receive message from broker.")
	s.mutex.Lock()
	defer s.mutex.Unlock()

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

		log.Println("WARN : topic id was not found for ", topic, ".")
		return
	}

	// send message
	s.Conn.WriteToUDP(m.Marshall(), s.Remote)
}

/**
 *
 */
func (s *TransportSnSession) sendMqttMessageLoop() {
	for {
		select {
		case msg := <-s.sendBuffer:
			switch msgi := msg.(type) {
			case *message.Publish:
				s.doPublish(msgi)
			case *message.Subscribe:
				s.doSubscribe(msgi)
			}
		case <-s.shutDown:
			return
		default:
			continue
		}
	}
}

/**
 *
 */
func (s *TransportSnSession) doPublish(m *message.Publish) {
	var topicName string
	var ok bool
	if message.TopicIdType(m.Flags) == message.MQTTSN_TIDT_NORMAL {
		// search topic name from topic id
		topicName, ok = s.LoadTopic(m.TopicId)
		if ok == false {
			// error handling
			log.Println("ERROR : topic was not found.")
			puback := message.NewPubAck(
				m.TopicId,
				m.MsgId,
				message.MQTTSN_RC_REJECTED_INVALID_TOPIC_ID)
			// send
			s.Conn.WriteToUDP(puback.Marshall(), s.Remote)

			return
		}
	} else if message.TopicIdType(m.Flags) == message.MQTTSN_TIDT_PREDEFINED {
		// search topic name from topic id
		topicName, ok = s.LoadPredefTopic(m.TopicId)
		if ok == false {
			// error handling
			log.Println("ERROR : topic was not found.")
			puback := message.NewPubAck(
				m.TopicId,
				m.MsgId,
				message.MQTTSN_RC_REJECTED_INVALID_TOPIC_ID)
			// send
			s.Conn.WriteToUDP(puback.Marshall(), s.Remote)

			return
		}
	} else if message.TopicIdType(m.Flags) == message.MQTTSN_TIDT_SHORT_NAME {
		topicName = m.TopicName
	} else {
		log.Println("ERROR : invalid TopicIdType ", message.TopicIdType(m.Flags))
		return
	}

	// get qos
	qos := message.QosFlag(m.Flags)

	// send Publish to Mqtt Broker.
	// Topic, Qos, Retain, Payload
	token := s.mqttClient.Publish(topicName, qos, false, m.Data)

	if qos == 1 {
		// if qos 1
		go waitPubAck(token, s.MqttSnSession, m.TopicId, m.MsgId)
	} else if qos == 2 {
		// elif qos 2
		// TODO: send PubRec
		// TODO: wait PubComp
	}
}

/**
 *
 */
func (s *TransportSnSession) doSubscribe(m *message.Subscribe) {
	// if topic include wildcard, set topicId as 0x0000
	// else regist topic to client-session instance and assign topiId
	var topicId uint16
	topicId = uint16(0x0000)

	// return code
	var rc byte = message.MQTTSN_RC_ACCEPTED

	switch message.TopicIdType(m.Flags) {
	// if TopicIdType is NORMAL, regist it
	case message.MQTTSN_TIDT_NORMAL:
		// check topic is wildcarded or not
		topics := strings.Split(m.TopicName, "/")
		if IsWildCarded(topics) != true {
			topicId = s.StoreTopic(m.TopicName)
		}

		// send subscribe to broker
		mqtt := s.mqttClient
		qos := message.QosFlag(m.Flags)
		mqtt.Subscribe(m.TopicName, qos, s.OnPublish)

	// else if PREDEFINED, get TopicName and Subscribe to Broker
	case message.MQTTSN_TIDT_PREDEFINED:
		// Get topicId
		topicId = m.TopicId

		// get topic name and subscribe to broker
		topicName, ok := s.LoadPredefTopic(topicId)
		if ok != true {
			// TODO: error handling
			rc = message.MQTTSN_RC_REJECTED_INVALID_TOPIC_ID
		}

		// send subscribe to broker
		mqtt := s.mqttClient
		qos := message.QosFlag(m.Flags)
		mqtt.Subscribe(topicName, qos, s.OnPublish)

	// else if SHORT_NAME, subscribe to broker
	case message.MQTTSN_TIDT_SHORT_NAME:
		// TODO: implement
	}

	// send subscribe to broker
	suback := message.NewSubAck(
		topicId,
		m.MsgId,
		rc)

	// send suback
	s.Conn.WriteToUDP(suback.Marshall(), s.Remote)
}
