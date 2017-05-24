package ceratopogon

import (
	"github.com/KMotiko/ceratopogon/messages"
	"log"
	"net"
	"strings"
	"sync"
)

func NewTransportGateway(config *GatewayConfig) *TransportGateway {
	g := &TransportGateway{
		sync.RWMutex{}, make(map[string]*TransportSnSession), config}
	return g
}

func (g *TransportGateway) HandlePacket(conn *net.UDPConn, remote *net.UDPAddr, packet []byte) {
	// parse message
	m := message.UnMarshall(packet)

	// handle message
	switch mi := m.(type) {
	case *message.MqttSnHeader:
		switch mi.MsgType {
		case message.MQTTSNT_WILLTOPICREQ:
			// WillTopicReq
			g.handleWillTopicReq(conn, remote, mi)
		case message.MQTTSNT_WILLMSGREQ:
			// WillMsgReq
			g.handleWillMsgReq(conn, remote, mi)
		case message.MQTTSNT_PINGREQ:
			// PingResp
			g.handlePingResp(conn, remote, mi)
		}
	case *message.Advertise:
		g.handleAdvertise(conn, remote, mi)
	case *message.SearchGw:
		g.handleSearchGw(conn, remote, mi)
	case *message.GwInfo:
		g.handleGwInfo(conn, remote, mi)
	case *message.Connect:
		g.handleConnect(conn, remote, mi)
	case *message.ConnAck:
		g.handleConnAck(conn, remote, mi)
	case *message.WillTopic:
		g.handleWillTopic(conn, remote, mi)
	case *message.WillMsg:
		g.handleWillMsg(conn, remote, mi)
	case *message.Register:
		g.handleRegister(conn, remote, mi)
	case *message.RegAck:
		g.handleRegAck(conn, remote, mi)
	case *message.Publish:
		g.handlePublish(conn, remote, mi)
	case *message.PubAck:
		g.handlePubAck(conn, remote, mi)
	case *message.PubRec:
		switch mi.Header.MsgType {
		case message.MQTTSNT_PUBREC:
			g.handlePubRec(conn, remote, mi)
		case message.MQTTSNT_PUBCOMP:
			//PubComp:
			g.handlePubComp(conn, remote, mi)
		case message.MQTTSNT_PUBREL:
			//PubRel:
			g.handlePubRel(conn, remote, mi)
		}
	case *message.Subscribe:
		switch mi.Header.MsgType {
		case message.MQTTSNT_SUBSCRIBE:
			// Subscribe
			g.handleSubscribe(conn, remote, mi)
		case message.MQTTSNT_UNSUBSCRIBE:
			// UnSubscribe
			g.handleUnSubscribe(conn, remote, mi)
		}
	case *message.SubAck:
		g.handleSubAck(conn, remote, mi)
	case *message.UnSubAck:
		g.handleUnSubAck(conn, remote, mi)
	case *message.PingReq:
		g.handlePingReq(conn, remote, mi)
	case *message.DisConnect:
		g.handleDisConnect(conn, remote, mi)
	case *message.WillTopicUpd:
		g.handleWillTopicUpd(conn, remote, mi)
	case *message.WillTopicResp:
		switch mi.Header.MsgType {
		case message.MQTTSNT_WILLTOPICRESP:
			// WillTopicResp
			g.handleWillTopicResp(conn, remote, mi)
		case message.MQTTSNT_WILLMSGRESP:
			// WillMsgResp
			g.handleWillMsgResp(conn, remote, mi)
		}
	case *message.WillMsgUpd:
		g.handleWillMsgUpd(conn, remote, mi)
	}
}

/*********************************************/
/* Advertise                                 */
/*********************************************/
func (g *TransportGateway) handleAdvertise(conn *net.UDPConn, remote *net.UDPAddr, m *message.Advertise) {
	// TODO: implement
}

/*********************************************/
/* SearchGw                                  */
/*********************************************/
func (g *TransportGateway) handleSearchGw(conn *net.UDPConn, remote *net.UDPAddr, m *message.SearchGw) {
	// TODO: implement
}

/*********************************************/
/* GwInfo                                    */
/*********************************************/
func (g *TransportGateway) handleGwInfo(conn *net.UDPConn, remote *net.UDPAddr, m *message.GwInfo) {
	// TODO: implement
}

/*********************************************/
/* Connect                                   */
/*********************************************/
func (g *TransportGateway) handleConnect(conn *net.UDPConn, remote *net.UDPAddr, m *message.Connect) {
	log.Println("handle Connect")
	// TODO: check connected client is aliready registerd or not

	// create mqtt-sn session instance
	s := NewTransportSnSession(m.ClientId, conn, remote)

	// connect to mqtt broker
	s.ConnectToBroker(
		g.Config.BrokerHost,
		g.Config.BrokerPort,
		g.Config.BrokerUser,
		g.Config.BrokerPassword)

	// send conn ack
	ack := message.NewConnAck()
	ack.ReturnCode = message.MQTTSN_RC_ACCEPTED
	packet := ack.Marshall()
	conn.WriteToUDP(packet, remote)

	// lock
	g.mutex.Lock()
	defer g.mutex.Unlock()

	// add session to map
	g.MqttSnSessions[remote.String()] = s
}

/*********************************************/
/* ConnAck                                   */
/*********************************************/
func (g *TransportGateway) handleConnAck(conn *net.UDPConn, remote *net.UDPAddr, m *message.ConnAck) {
	// TODO: implement
}

/*********************************************/
/* WillTopicReq                              */
/*********************************************/
func (g *TransportGateway) handleWillTopicReq(conn *net.UDPConn, remote *net.UDPAddr, m *message.MqttSnHeader) {
	// TODO: implement
}

/*********************************************/
/* WillTopic                                 */
/*********************************************/
func (g *TransportGateway) handleWillTopic(conn *net.UDPConn, remote *net.UDPAddr, m *message.WillTopic) {
	// TODO: implement
}

/*********************************************/
/* WillMsgReq                                */
/*********************************************/
func (g *TransportGateway) handleWillMsgReq(conn *net.UDPConn, remote *net.UDPAddr, m *message.MqttSnHeader) {
	// TODO: implement
}

/*********************************************/
/* WillMsg                                   */
/*********************************************/
func (g *TransportGateway) handleWillMsg(conn *net.UDPConn, remote *net.UDPAddr, m *message.WillMsg) {
	// TODO: implement
}

/*********************************************/
/* Register                                  */
/*********************************************/
func (g *TransportGateway) handleRegister(conn *net.UDPConn, remote *net.UDPAddr, m *message.Register) {
	log.Println("handle Register")

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
func (g *TransportGateway) handleRegAck(conn *net.UDPConn, remote *net.UDPAddr, m *message.RegAck) {
	log.Println("handle RegAck")
	if m.ReturnCode != message.MQTTSN_RC_ACCEPTED {
		log.Println("ERROR : RegAck not accepted. Reqson = ", m.ReturnCode)
	}
}

/*********************************************/
/* Publish                                   */
/*********************************************/
func (g *TransportGateway) handlePublish(conn *net.UDPConn, remote *net.UDPAddr, m *message.Publish) {
	log.Println("handle Publish")

	// get mqttsn session
	s, ok := g.MqttSnSessions[remote.String()]
	if ok == false {
		// TODO: error handling
	}

	var topicName string
	if message.TopicIdType(m.Flags) == message.MQTTSN_TIDT_PREDEFINED {
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
	} else if message.TopicIdType(m.Flags) == message.MQTTSN_TIDT_SHORT_NAME {
		topicName = m.TopicName
	} else {
		log.Println("ERROR : invalid TopicIdType ", message.TopicIdType(m.Flags))
	}

	// send publish to broker
	mqtt := s.mqttClient

	// get qos
	qos := message.QosFlag(m.Flags)

	// send Publish to Mqtt Broker.
	// Topic, Qos, Retain, Payload
	mqtt.Publish(topicName, qos, false, m.Data)

	if qos == 0 {
		// if qos 0
		// nothing to do
	} else if qos == 1 {
		// elif qos 1
		// Create PubAck
		puback := message.NewPubAck(m.TopicId, m.MsgId, message.MQTTSN_RC_ACCEPTED)

		// send
		conn.WriteToUDP(puback.Marshall(), remote)
	} else if qos == 2 {
		// elif qos 2
		// send PubRec
	}
}

/*********************************************/
/* PubAck                                    */
/*********************************************/
func (g *TransportGateway) handlePubAck(conn *net.UDPConn, remote *net.UDPAddr, m *message.PubAck) {
	// TODO: implement
}

/*********************************************/
/* PubRec                                    */
/*********************************************/
func (g *TransportGateway) handlePubRec(conn *net.UDPConn, remote *net.UDPAddr, m *message.PubRec) {
	// TODO: implement
}

/*********************************************/
/* PubRel                                    */
/*********************************************/
func (g *TransportGateway) handlePubRel(conn *net.UDPConn, remote *net.UDPAddr, m *message.PubRec) {
	// TODO: implement
}

/*********************************************/
/* PubComp                                   */
/*********************************************/
func (g *TransportGateway) handlePubComp(conn *net.UDPConn, remote *net.UDPAddr, m *message.PubRec) {
	// TODO: implement
}

/*********************************************/
/* Subscribe                                 */
/*********************************************/
func (g *TransportGateway) handleSubscribe(conn *net.UDPConn, remote *net.UDPAddr, m *message.Subscribe) {
	log.Println("handle Subscribe")

	// TODO: add lock?
	// get mqttsn session
	s, ok := g.MqttSnSessions[remote.String()]
	if ok == false {
		// TODO: error handling
	}

	// if topic include wildcard, set topicId as 0x0000
	// else regist topic to client-session instance and assign topiId
	var topicId uint16
	topicId = uint16(0x0000)

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
		// get topic name and subscribe to broker
		topicName, ok := s.Topics.LoadTopic(m.TopicId)
		if ok != true {
			// TODO: error handling
		}

		// Get topicId
		topicId = m.TopicId

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
		message.MQTTSN_RC_ACCEPTED)

	// send suback
	conn.WriteToUDP(suback.Marshall(), remote)
}

/*********************************************/
/* SubAck                                 */
/*********************************************/
func (g *TransportGateway) handleSubAck(conn *net.UDPConn, remote *net.UDPAddr, m *message.SubAck) {
	// TODO: implement
}

/*********************************************/
/* UnSubscribe                               */
/*********************************************/
func (g *TransportGateway) handleUnSubscribe(conn *net.UDPConn, remote *net.UDPAddr, m *message.Subscribe) {
	// TODO: implement
}

/*********************************************/
/* UnSubAck                                  */
/*********************************************/
func (g *TransportGateway) handleUnSubAck(conn *net.UDPConn, remote *net.UDPAddr, m *message.UnSubAck) {
	// TODO: implement
}

/*********************************************/
/* PingReq                                  */
/*********************************************/
func (g *TransportGateway) handlePingReq(conn *net.UDPConn, remote *net.UDPAddr, m *message.PingReq) {
	// TODO: implement
}

/*********************************************/
/* PingResp                                  */
/*********************************************/
func (g *TransportGateway) handlePingResp(conn *net.UDPConn, remote *net.UDPAddr, m *message.MqttSnHeader) {
	// TODO: implement
}

/*********************************************/
/* DisConnect                                */
/*********************************************/
func (g *TransportGateway) handleDisConnect(conn *net.UDPConn, remote *net.UDPAddr, m *message.DisConnect) {
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

	// TODO: disconnect mqtt connection with broker ?

	// send DisConnect
	dc := message.NewDisConnect(0)
	packet := dc.Marshall()
	conn.WriteToUDP(packet, remote)
}

/*********************************************/
/* WillTopicUpd                              */
/*********************************************/
func (g *TransportGateway) handleWillTopicUpd(conn *net.UDPConn, remote *net.UDPAddr, m *message.WillTopicUpd) {
	// TODO: implement
}

/*********************************************/
/* WillMsgUpd                                */
/*********************************************/
func (g *TransportGateway) handleWillMsgUpd(conn *net.UDPConn, remote *net.UDPAddr, m *message.WillMsgUpd) {
	// TODO: implement
}

/*********************************************/
/* WillTopicResp                             */
/*********************************************/
func (g *TransportGateway) handleWillTopicResp(conn *net.UDPConn, remote *net.UDPAddr, m *message.WillTopicResp) {
	// TODO: implement
}

/*********************************************/
/* WillMsgResp                             */
/*********************************************/
func (g *TransportGateway) handleWillMsgResp(conn *net.UDPConn, remote *net.UDPAddr, m *message.WillTopicResp) {
	// TODO: implement
}
