package ceratopoGon

import (
	"github.com/Kmotiko/ceratopoGon/env"
	"github.com/Kmotiko/ceratopoGon/messages"
	"log"
	"net"
)

func NewTransportGateway(
	config *GatewayConfig,
	predefTopics PredefinedTopics) *TransportGateway {
	g := &TransportGateway{
		make(map[string]*TransportSnSession),
		config,
		predefTopics}
	return g
}

func (g *TransportGateway) StartUp() error {
	// launch server loop
	return serverLoop(g, g.Config.Host, g.Config.Port)
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
	log.Println("Connected from ", m.ClientId, " : ", remote.String())

	var s *TransportSnSession = nil

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
			s.shutDown <- true
		}

		// create new session
		s = NewTransportSnSession(
			m.ClientId,
			conn,
			remote,
			g.Config.BrokerHost,
			g.Config.BrokerPort,
			g.Config.BrokerUser,
			g.Config.BrokerPassword)

		// read predefined topics
		if topics, ok := g.predefTopics[m.ClientId]; ok {
			for key, value := range topics {
				s.StorePredefTopicWithId(key, value)
			}
		}

		// connect to mqtt broker
		s.ConnectToBroker(true)

		go s.sendMqttMessageLoop()
	}

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
		log.Println("ERROR : MqttSn session not found for remote", remote.String())
		return
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
	if env.DEBUG {
		log.Println("handle Publish")
		log.Println("Published from : ", remote.String(), ", TopicID : ", m.TopicId)
	}

	// get mqttsn session
	s, ok := g.MqttSnSessions[remote.String()]
	if ok == false {
		log.Println("ERROR : MqttSn session not found for remote", remote.String())
		return
	}
	s.sendBuffer <- m

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

	// get mqttsn session
	s, ok := g.MqttSnSessions[remote.String()]
	if ok == false {
		log.Println("ERROR : MqttSn session not found for remote", remote.String())
		return
	}

	s.sendBuffer <- m
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
	log.Println("DisConnect : ", remote.String())

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
