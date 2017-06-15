package ceratopoGon

import (
	"encoding/json"
	"io/ioutil"
)

type GatewayConfig struct {
	IsAggregate    bool
	Host           string
	Port           int
	BrokerHost     string
	BrokerPort     int
	BrokerUser     string
	BrokerPassword string
	LogFilePath    string
}

func ParseConfig(filename string) (*GatewayConfig, error) {
	config := new(GatewayConfig)
	jsonStr, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonStr, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
