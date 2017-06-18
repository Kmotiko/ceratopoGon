package ceratopoGon

import (
	"github.com/Kmotiko/ceratopoGon/messages"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"net"
	"strconv"
	"strings"
)

func NewAggregatingGateway(config *GatewayConfig,
	predefTopics PredefinedTopics) *AggregatingGateway {
	g := &AggregatingGateway{
		MqttSnSessions: make(map[string]*MqttSnSession),
		Config:         config,
		predefTopics:   predefTopics}
	return g
}

func (g *AggregatingGateway) StartUp() error {
	// create opts
	addr := "tcp://" + g.Config.BrokerHost + ":" + strconv.Itoa(g.Config.BrokerPort)
	opts := MQTT.NewClientOptions().AddBroker(addr)
	opts.SetClientID("ceratopogon")
	opts.SetUsername(g.Config.BrokerUser)
	opts.SetPassword(g.Config.BrokerPassword)

	// create client instance
	// TODO: add lock
	g.mqttClient = (MQTT.NewClient(opts))

	// connect
	if token := g.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
		return token.Error()
	}

	// launch server loop
	return serverLoop(g, g.Config.Host, g.Config.Port)
}

func (g *AggregatingGateway) HandlePacket(conn *net.UDPConn, remote *net.UDPAddr, packet []byte) {
	// parse message
	m := message.UnMarshall(packet)

	// handle message
	switch mi := m.(type) {
	case *message.MqttSnHeader:
		switch mi.MsgType {
		case message.MQTTSNT_WILLTOPICREQ:
			// WillTopicReq
			go g.handleWillTopicReq(conn, remote, mi)
		case message.MQTTSNT_WILLMSGREQ:
			// WillMsgReq
			go g.handleWillMsgReq(conn, remote, mi)
		case message.MQTTSNT_PINGREQ:
			// PingResp
			go g.handlePingResp(conn, remote, mi)
		}
	case *message.Advertise:
		go g.handleAdvertise(conn, remote, mi)
	case *message.SearchGw:
		go g.handleSearchGw(conn, remote, mi)
	case *message.GwInfo:
		go g.handleGwInfo(conn, remote, mi)
	case *message.Connect:
		go g.handleConnect(conn, remote, mi)
	case *message.ConnAck:
		go g.handleConnAck(conn, remote, mi)
	case *message.WillTopic:
		go g.handleWillTopic(conn, remote, mi)
	case *message.WillMsg:
		go g.handleWillMsg(conn, remote, mi)
	case *message.Register:
		go g.handleRegister(conn, remote, mi)
	case *message.RegAck:
		go g.handleRegAck(conn, remote, mi)
	case *message.Publish:
		go g.handlePublish(conn, remote, mi)
	case *message.PubAck:
		go g.handlePubAck(conn, remote, mi)
	case *message.PubRec:
		switch mi.Header.MsgType {
		case message.MQTTSNT_PUBREC:
			go g.handlePubRec(conn, remote, mi)
		case message.MQTTSNT_PUBCOMP:
			//PubComp:
			go g.handlePubComp(conn, remote, mi)
		case message.MQTTSNT_PUBREL:
			//PubRel:
			go g.handlePubRel(conn, remote, mi)
		}
	case *message.Subscribe:
		switch mi.Header.MsgType {
		case message.MQTTSNT_SUBSCRIBE:
			// Subscribe
			go g.handleSubscribe(conn, remote, mi)
		case message.MQTTSNT_UNSUBSCRIBE:
			// UnSubscribe
			go g.handleUnSubscribe(conn, remote, mi)
		}
	case *message.SubAck:
		go g.handleSubAck(conn, remote, mi)
	case *message.UnSubAck:
		go g.handleUnSubAck(conn, remote, mi)
	case *message.PingReq:
		go g.handlePingReq(conn, remote, mi)
	case *message.DisConnect:
		go g.handleDisConnect(conn, remote, mi)
	case *message.WillTopicUpd:
		go g.handleWillTopicUpd(conn, remote, mi)
	case *message.WillTopicResp:
		switch mi.Header.MsgType {
		case message.MQTTSNT_WILLTOPICRESP:
			// WillTopicResp
			go g.handleWillTopicResp(conn, remote, mi)
		case message.MQTTSNT_WILLMSGRESP:
			// WillMsgResp
			go g.handleWillMsgResp(conn, remote, mi)
		}
	case *message.WillMsgUpd:
		go g.handleWillMsgUpd(conn, remote, mi)
	}
}

/*********************************************/
/* Advertise                                 */
/*********************************************/
func (g *AggregatingGateway) handleAdvertise(conn *net.UDPConn, remote *net.UDPAddr, m *message.Advertise) {
	// TODO: implement
}

/*********************************************/
/* SearchGw                                  */
/*********************************************/
func (g *AggregatingGateway) handleSearchGw(conn *net.UDPConn, remote *net.UDPAddr, m *message.SearchGw) {
	// TODO: implement
}

/*********************************************/
/* GwInfo                                    */
/*********************************************/
func (g *AggregatingGateway) handleGwInfo(conn *net.UDPConn, remote *net.UDPAddr, m *message.GwInfo) {
	// TODO: implement
}

/*********************************************/
/* Connect                                   */
/*********************************************/
func (g *AggregatingGateway) handleConnect(conn *net.UDPConn, remote *net.UDPAddr, m *message.Connect) {
	log.Println("handle Connect")
	// lock
	g.mutex.Lock()
	defer g.mutex.Unlock()

	var s *MqttSnSession = nil

	// check connected client is already registerd or not
	// and cleansession is required or not
	for _, val := range g.MqttSnSessions {
		if val.ClientId == m.ClientId {
			s = val
		}
	}

	// if client is already registerd and cleansession is false
	// reuse mqtt-sn session
	if s != nil && !message.HasCleanSessionFlag(m.Flags) {
		delete(g.MqttSnSessions, s.Remote.String())
		s.Remote = remote
		s.Conn = conn
	} else {
		// otherwise(cleansession is true or first time to connect)
		// create new mqtt-sn session instance
		if s != nil {
			delete(g.MqttSnSessions, s.Remote.String())
		}

		// create new session
		s = NewMqttSnSession(m.ClientId, conn, remote)

		// read predefined topics
		if topics, ok := g.predefTopics[m.ClientId]; ok {
			for key, value := range topics {
				s.StorePredefTopicWithId(key, value)
			}
		}
	}

	// TODO: check will flags
	// TODO: support WILLTOPICREQ, WILLMSGREQ if will flag is true

	// add session to map
	g.MqttSnSessions[remote.String()] = s

	// send conn ack
	ack := message.NewConnAck()
	ack.ReturnCode = message.MQTTSN_RC_ACCEPTED
	packet := ack.Marshall()
	conn.WriteToUDP(packet, remote)
}

/*********************************************/
/* ConnAck                                   */
/*********************************************/
func (g *AggregatingGateway) handleConnAck(conn *net.UDPConn, remote *net.UDPAddr, m *message.ConnAck) {
	// TODO: implement
}

/*********************************************/
/* WillTopicReq                              */
/*********************************************/
func (g *AggregatingGateway) handleWillTopicReq(conn *net.UDPConn, remote *net.UDPAddr, m *message.MqttSnHeader) {
	// TODO: implement
}

/*********************************************/
/* WillTopic                                 */
/*********************************************/
func (g *AggregatingGateway) handleWillTopic(conn *net.UDPConn, remote *net.UDPAddr, m *message.WillTopic) {
	// TODO: implement
}

/*********************************************/
/* WillMsgReq                                */
/*********************************************/
func (g *AggregatingGateway) handleWillMsgReq(conn *net.UDPConn, remote *net.UDPAddr, m *message.MqttSnHeader) {
	// TODO: implement
}

/*********************************************/
/* WillMsg                                   */
/*********************************************/
func (g *AggregatingGateway) handleWillMsg(conn *net.UDPConn, remote *net.UDPAddr, m *message.WillMsg) {
	// TODO: implement
}

/*********************************************/
/* Register                                  */
/*********************************************/
func (g *AggregatingGateway) handleRegister(conn *net.UDPConn, remote *net.UDPAddr, m *message.Register) {
	log.Println("handle Register")
	// lock
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// get mqttsn session
	s, ok := g.MqttSnSessions[remote.String()]
	if ok == false {
		// TODO: error handling
	}

	// search topicid
	topicId, ok := s.LoadTopicId(m.TopicName)
	if !ok {
		// store topic to map
		topicId = s.StoreTopic(m.TopicName)
	}

	// Create RegAck with topicid
	regAck := message.NewRegAck(
		topicId, m.MsgId, message.MQTTSN_RC_ACCEPTED)

	// send RegAck
	conn.WriteToUDP(regAck.Marshall(), remote)
}

/*********************************************/
/* RegAck                                    */
/*********************************************/
func (g *AggregatingGateway) handleRegAck(conn *net.UDPConn, remote *net.UDPAddr, m *message.RegAck) {
	log.Println("handle RegAck")
	if m.ReturnCode != message.MQTTSN_RC_ACCEPTED {
		log.Println("ERROR : RegAck not accepted. Reqson = ", m.ReturnCode)
	}
}

/*********************************************/
/* Publish                                   */
/*********************************************/
func (g *AggregatingGateway) handlePublish(conn *net.UDPConn, remote *net.UDPAddr, m *message.Publish) {
	log.Println("handle Publish")
	// lock
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// get mqttsn session
	s, ok := g.MqttSnSessions[remote.String()]
	if ok == false {
		log.Println("ERROR : MqttSn session not found for remote", remote.String())
		return
	}

	var topicName string
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
			conn.WriteToUDP(puback.Marshall(), remote)

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
			conn.WriteToUDP(puback.Marshall(), remote)

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
	// TODO: add lock ?
	token := g.mqttClient.Publish(topicName, qos, false, m.Data)

	if qos == 1 {
		// if qos 1
		go waitPubAck(token, s, m.TopicId, m.MsgId)
	} else if qos == 2 {
		// elif qos 2
		// TODO: send PubRec
		// TODO: wait PubComp
	}
}

/*********************************************/
/* PubAck                                    */
/*********************************************/
func (g *AggregatingGateway) handlePubAck(conn *net.UDPConn, remote *net.UDPAddr, m *message.PubAck) {
	// TODO: implement
}

/*********************************************/
/* PubRec                                    */
/*********************************************/
func (g *AggregatingGateway) handlePubRec(conn *net.UDPConn, remote *net.UDPAddr, m *message.PubRec) {
	// TODO: implement
}

/*********************************************/
/* PubRel                                    */
/*********************************************/
func (g *AggregatingGateway) handlePubRel(conn *net.UDPConn, remote *net.UDPAddr, m *message.PubRec) {
	// TODO: implement
}

/*********************************************/
/* PubComp                                   */
/*********************************************/
func (g *AggregatingGateway) handlePubComp(conn *net.UDPConn, remote *net.UDPAddr, m *message.PubRec) {
	// TODO: implement
}

/*********************************************/
/* Subscribe                                 */
/*********************************************/
func (g *AggregatingGateway) handleSubscribe(conn *net.UDPConn, remote *net.UDPAddr, m *message.Subscribe) {
	log.Println("handle Subscribe")

	// get mqttsn session
	g.mutex.Lock()
	s, ok := g.MqttSnSessions[remote.String()]
	if ok == false {
		log.Println("ERROR : MqttSn session not found for remote", remote.String())
		g.mutex.Unlock()
		return
	}
	g.mutex.Unlock()

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
		subscribers := GetTopicEntry().AppendSubscriber(
			m.TopicName,
			s,
			message.MQTTSN_TIDT_NORMAL)

		// if first subscribers, send subscribe to broker
		if len(subscribers) == 1 {
			qos := message.QosFlag(m.Flags)
			if token := g.mqttClient.Subscribe(
				(m.TopicName), byte(qos), g.OnPublish); token.Wait() && token.Error() != nil {
				// TODO: error handling
				log.Println("failed to send Subscribe to broker : ", token.Error())
			}
		}

	// else if PREDEFINED, get TopicName and Subscribe to Broker
	case message.MQTTSN_TIDT_PREDEFINED:
		// get topic name and subscribe to broker
		topicName, ok := s.Topics.LoadTopic(m.TopicId)
		if ok != true {
			// TODO: error handling
			rc = message.MQTTSN_RC_REJECTED_INVALID_TOPIC_ID
		}

		// Get topicId
		topicId = m.TopicId

		// PREDEFINED topic will not be wildcarded
		subscribers := GetTopicEntry().AppendSubscriber(
			topicName,
			s,
			message.MQTTSN_TIDT_PREDEFINED)

		// if first subscribers, send subscribe to broker
		if len(subscribers) == 1 {
			qos := message.QosFlag(m.Flags)
			if token := g.mqttClient.Subscribe(
				(topicName), byte(qos), g.OnPublish); token.Wait() && token.Error() != nil {
				// TODO: error handling
				log.Println("failed to send Subscribe to broker : ", token.Error())
			}
		}

	// else if SHORT_NAME, subscribe to broker
	case message.MQTTSN_TIDT_SHORT_NAME:
		// TODO: implement

		// append to TopicNodes
	}

	// send suback to client
	suback := message.NewSubAck(
		topicId,
		m.MsgId,
		rc)

	// send suback
	conn.WriteToUDP(suback.Marshall(), remote)
}

/*********************************************/
/* SubAck                                 */
/*********************************************/
func (g *AggregatingGateway) handleSubAck(conn *net.UDPConn, remote *net.UDPAddr, m *message.SubAck) {
	// TODO: implement
}

/*********************************************/
/* UnSubscribe                               */
/*********************************************/
func (g *AggregatingGateway) handleUnSubscribe(conn *net.UDPConn, remote *net.UDPAddr, m *message.Subscribe) {
	// TODO: implement
}

/*********************************************/
/* UnSubAck                                  */
/*********************************************/
func (g *AggregatingGateway) handleUnSubAck(conn *net.UDPConn, remote *net.UDPAddr, m *message.UnSubAck) {
	// TODO: implement
}

/*********************************************/
/* PingReq                                  */
/*********************************************/
func (g *AggregatingGateway) handlePingReq(conn *net.UDPConn, remote *net.UDPAddr, m *message.PingReq) {
	// TODO: implement
}

/*********************************************/
/* PingResp                                  */
/*********************************************/
func (g *AggregatingGateway) handlePingResp(conn *net.UDPConn, remote *net.UDPAddr, m *message.MqttSnHeader) {
	// TODO: implement
}

/*********************************************/
/* DisConnect                                */
/*********************************************/
func (g *AggregatingGateway) handleDisConnect(conn *net.UDPConn, remote *net.UDPAddr, m *message.DisConnect) {
	log.Println("handle DisConnect")

	// lock
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// TODO: implement
	if m.Duration > 0 {
		// TODO: set session state to SLEEP
	} else {
		// TODO: set session state to DISCONNECT
	}

	// send DisConnect
	dc := message.NewDisConnect(0)
	packet := dc.Marshall()
	conn.WriteToUDP(packet, remote)
}

/*********************************************/
/* WillTopicUpd                              */
/*********************************************/
func (g *AggregatingGateway) handleWillTopicUpd(conn *net.UDPConn, remote *net.UDPAddr, m *message.WillTopicUpd) {
	// TODO: implement
}

/*********************************************/
/* WillMsgUpd                                */
/*********************************************/
func (g *AggregatingGateway) handleWillMsgUpd(conn *net.UDPConn, remote *net.UDPAddr, m *message.WillMsgUpd) {
	// TODO: implement
}

/*********************************************/
/* WillTopicResp                             */
/*********************************************/
func (g *AggregatingGateway) handleWillTopicResp(conn *net.UDPConn, remote *net.UDPAddr, m *message.WillTopicResp) {
	// TODO: implement
}

/*********************************************/
/* WillMsgResp                               */
/*********************************************/
func (g *AggregatingGateway) handleWillMsgResp(conn *net.UDPConn, remote *net.UDPAddr, m *message.WillTopicResp) {
	// TODO: implement
}

/*********************************************/
/* OnPublish                                 */
/*********************************************/
// handle message from broker
func (g *AggregatingGateway) OnPublish(client MQTT.Client, msg MQTT.Message) {
	log.Println("on publish. Receive message from broker.")
	// get subscribers
	topic := msg.Topic()
	subscribers := GetTopicEntry().GetSubscriber(topic)

	// for each subscriber
	for _, subscriber := range subscribers {
		// TODO: check subscriber(MqttSnSession)'s state.
		// if subscriber is sleep, gateway must buffer the message.

		var m *message.Publish
		session := subscriber.Session

		// get TopicID Type
		tidType := subscriber.TidType
		switch tidType {
		case message.MQTTSN_TIDT_NORMAL:
			// get topicid
			topicId, ok := session.LoadTopicId(topic)

			// if not found
			if !ok {
				// TODO: implement
				log.Println("[Info] Normal topic id was not found for ", topic, ".")

				// TODO: implement wildcarded route
				continue

			} else {
				// process as fixed topicid
				// qos, retain, topicId, msgId, data
				msgId := session.NextMsgId()
				m = message.NewPublishNormal(
					1, false, topicId, msgId, msg.Payload())
			}

		case message.MQTTSN_TIDT_PREDEFINED:
			// get topicid
			topicId, ok := session.LoadPredefTopicId(topic)

			// if not found
			if !ok {
				// PREDEFINED TOPIC is must be fixed
				log.Println("[Error] Predefined topic id was not found for ", topic, ".")
				continue

			} else {
				// process as fixed topicid
				// qos, retain, topicId, msgId, data
				msgId := session.NextMsgId()
				m = message.NewPublishPredefined(
					1, false, topicId, msgId, msg.Payload())
			}

		case message.MQTTSN_TIDT_SHORT_NAME:

		}

		// send message
		session.Conn.WriteToUDP(m.Marshall(), session.Remote)
	}
}
