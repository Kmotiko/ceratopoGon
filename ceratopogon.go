package main

import (
	"errors"
	"flag"
	"github.com/Kmotiko/ceratopoGon/gateway"
	"log"
	"os"
	"os/signal"
	"syscall"
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
	if topicFile != "" {
		err = ceratopoGon.InitPredefinedTopic(topicFile)
		if err != nil {
			return nil, err
		}
	}

	// initialize logger
	err = ceratopoGon.InitLogger(config.LogFilePath)
	if err != nil {
		return nil, err
	}

	// create signal chan
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGKILL)

	// create Gateway
	var gateway ceratopoGon.Gateway
	if config.IsAggregate {
		gateway = ceratopoGon.NewAggregatingGateway(config, signalChan)
	} else {
		gateway = ceratopoGon.NewTransparentGateway(config, signalChan)
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
		log.Println(errors.New("ERROR : failed to StartUp gateway"))
	}
	return
}
