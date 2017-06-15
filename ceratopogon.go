package main

import (
	"errors"
	"flag"
	"github.com/Kmotiko/ceratopoGon/gateway"
	"log"
)

func initialize() (ceratopoGon.Gateway, error) {
	var confFile string
	var topicFile string
	flag.StringVar(&confFile, "c", "ceratopogon.conf", "config file path")
	flag.StringVar(&topicFile, "t", "", "predefined topic file path")
	flag.Parse()

	// parse config
	config, err := ceratopoGon.ParseConfig(confFile)
	if err != nil {
		return nil, err
	}

	// parse topic file
	var topics ceratopoGon.PredefinedTopics
	if topicFile != "" {
		topics, err = ceratopoGon.LoadPredefinedTopics(topicFile)
		if err != nil {
			return nil, err
		}
	}

	// initialize logger
	err = ceratopoGon.InitLogger(config.LogFilePath)
	if err != nil {
		return nil, err
	}

	// create Gateway
	var gateway ceratopoGon.Gateway
	if config.IsAggregate {
		gateway = ceratopoGon.NewAggregatingGateway(config, topics)
	} else {
		gateway = ceratopoGon.NewTransportGateway(config, topics)
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
