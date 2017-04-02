package ceratopogon

import (
	// "github.com/KMotiko/ceratopogon/messages"
	"net"
)

type Gateway interface {
	HandlePacket(*net.UDPConn, *net.UDPAddr, []byte)
}

type AggregatingGateway struct {
	// sendBuffer chan *message.MqttSnMessage
	// mqttSnSessions map[string]MqttSnSession
	// mqttClient MqttClient
	// topics
	// config GatewayConfig
}

type TransportGateway struct {
	MqttSnSessions map[string]*TransportSnSession
	Config         *GatewayConfig
}
