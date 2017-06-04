package main

import (
	"errors"
	"flag"
	"github.com/KMotiko/ceratopogon/gateway"
	"log"
)

func initialize() (ceratopogon.Gateway, error) {
	var confFile string
	var topicFile string
	flag.StringVar(&confFile, "c", "ceratopogon.conf", "config file path")
	flag.StringVar(&topicFile, "t", "", "predefined topic file path")
	flag.Parse()

	// parse config
	config, err := ceratopogon.ParseConfig(confFile)
	if err != nil {
		return nil, err
	}

	// parse topic file
	var topics ceratopogon.PredefinedTopics
	if topicFile != "" {
		topics, err = ceratopogon.LoadPredefinedTopics(topicFile)
		if err != nil {
			return nil, err
		}
	}

	// initialize logger
	err = ceratopogon.InitLogger(config.LogFilePath)
	if err != nil {
		return nil, err
	}

	// create Gateway
	var gateway ceratopogon.Gateway
	if config.IsAggregate {
		gateway = ceratopogon.NewAggregatingGateway(config, topics)
	} else {
		gateway = ceratopogon.NewTransportGateway(config, topics)
	}

	return gateway, nil
}

func main() {
	// initialize
	gateway, err := initialize()
	if err != nil {
		log.Println(err)
		return
	}

	// start server
	err = gateway.StartUp()
	if err != nil {
		log.Println(errors.New("failed to StartUp gateway"))
	}
	return
}
