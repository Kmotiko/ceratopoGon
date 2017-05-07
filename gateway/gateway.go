package ceratopogon

import (
	//	"github.com/KMotiko/ceratopogon/messages"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"net"
	"sync"
)

type Gateway interface {
	HandlePacket(*net.UDPConn, *net.UDPAddr, []byte)
}

type AggregatingGateway struct {
	mutex          sync.RWMutex
	MqttSnSessions map[string]*MqttSnSession
	Config         *GatewayConfig
	mqttClient     MQTT.Client
	// sendBuffer     chan *message.MqttSnMessage
	// topics Topic
}

type TransportGateway struct {
	mutex          sync.RWMutex
	MqttSnSessions map[string]*TransportSnSession
	Config         *GatewayConfig
}
