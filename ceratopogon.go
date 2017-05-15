package main

import (
	"flag"
	"github.com/KMotiko/ceratopogon/gateway"
	"log"
	"net"
	"strconv"
)

func serverLoop(config *ceratopogon.GatewayConfig) error {
	serverStr := config.Host + ":" + strconv.Itoa(config.Port)
	udpAddr, err := net.ResolveUDPAddr("udp", serverStr)
	conn, err := net.ListenUDP("udp", udpAddr)

	if err != nil {
		return err
	}

	buf := make([]byte, 2048)

	var gateway ceratopogon.Gateway
	if config.IsAggregate {
		gateway = ceratopogon.NewAggregatingGateway(config)
		gateway.(*(ceratopogon.AggregatingGateway)).StartUp()
	} else {
		gateway = ceratopogon.NewTransportGateway(config)
	}

	for {
		_, remote, err := conn.ReadFromUDP(buf)
		if err != nil {
			return err
		}

		// launch gateway
		go gateway.HandlePacket(conn, remote, buf)
	}
}

func parseConfig() (*ceratopogon.GatewayConfig, error) {
	var confFile string
	flag.StringVar(&confFile, "c", "ceratopogon.conf", "config file path")
	flag.Parse()

	config, err := ceratopogon.ParseConfig(confFile)
	if err != nil {
		return nil, err
	}

	err = ceratopogon.InitLogger(config.LogFilePath)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func main() {
	// initialize
	config, err := parseConfig()
	if err != nil {
		log.Println(err)
		return
	}

	// start server
	serverLoop(config)
}
