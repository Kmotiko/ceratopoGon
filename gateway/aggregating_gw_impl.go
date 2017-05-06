package ceratopogon

import (
	"github.com/KMotiko/ceratopogon/messages"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"net"
	"strconv"
)

func NewAggregatingGateway(config *GatewayConfig) *AggregatingGateway {
	g := &AggregatingGateway{
		MqttSnSessions: make(map[string]*MqttSnSession),
		Config:         config}
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
		return token.Error()
	}

	return nil
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
	// TODO: check connected client is already registerd or not

	// create mqtt-sn session instance
	s := NewMqttSnSession(m.ClientId, conn, remote)

	// send conn ack
	ack := message.NewConnAck()
	ack.ReturnCode = message.MQTTSN_RC_ACCEPTED
	packet := ack.Marshall()
	conn.WriteToUDP(packet, remote)

	// add session to map
	g.MqttSnSessions[remote.String()] = s
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

	// get qos
	qos := message.QosFlag(m.Flags)

	// send Publish to Mqtt Broker.
	// Topic, Qos, Retain, Payload
	// TODO: add lock
	g.mqttClient.Publish(topicName, qos, false, m.Data)

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

		// TODO: implement
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
	// TODO: implement

	// regist topic to gateway instance

	// if topic include wildcard, set topicId as 0x0000

	// send subscribe to broker

	// send suback
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
	// TODO: implement
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
/* WillMsgResp                             */
/*********************************************/
func (g *AggregatingGateway) handleWillMsgResp(conn *net.UDPConn, remote *net.UDPAddr, m *message.WillTopicResp) {
	// TODO: implement
}
