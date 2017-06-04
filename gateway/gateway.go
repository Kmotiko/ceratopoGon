package ceratopogon

import (
	"github.com/KMotiko/ceratopogon/messages"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"net"
	"strconv"
	"sync"
)

type Gateway interface {
	HandlePacket(*net.UDPConn, *net.UDPAddr, []byte)
	StartUp() error
}

type AggregatingGateway struct {
	mutex          sync.RWMutex
	MqttSnSessions map[string]*MqttSnSession
	Config         *GatewayConfig
	mqttClient     MQTT.Client
	predefTopics   PredefinedTopics
	// sendBuffer     chan *message.MqttSnMessage
	// topics Topic
}

type TransportGateway struct {
	mutex          sync.RWMutex
	MqttSnSessions map[string]*TransportSnSession
	Config         *GatewayConfig
	predefTopics   PredefinedTopics
}

func serverLoop(gateway Gateway, host string, port int) error {
	serverStr := host + ":" + strconv.Itoa(port)
	udpAddr, err := net.ResolveUDPAddr("udp", serverStr)
	conn, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		return err
	}

	buf := make([]byte, 64*1024)
	for {
		_, remote, err := conn.ReadFromUDP(buf)
		if err != nil {
			return err
		}

		// launch gateway
		go gateway.HandlePacket(conn, remote, buf)
	}
}

func waitPubAck(token MQTT.Token, s *MqttSnSession, topicId uint16, msgId uint16) {
	token.Wait()
	puback := message.NewPubAck(topicId, msgId, message.MQTTSN_RC_ACCEPTED)
	s.Conn.WriteToUDP(puback.Marshall(), s.Remote)
}

func waitPubComp(token MQTT.Token, s *MqttSnSession, topicId uint16, msgId uint16) {
	// TODO: implement
	// token.Wait()
	// pubcomp := message.NewPubComp()
	// s.Conn.WriteToUDP(pubcomp.Marshall(), s.Remote)
}
