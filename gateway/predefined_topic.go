package ceratopogon

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type PredefinedTopicNode struct {
	Name string `yaml:"Name"`
	Id   uint16 `yaml:"Id"`
}

type PredefinedTopics struct {
	ClientId string                `yaml:"ClientId"`
	Topics   []PredefinedTopicNode `yaml:"Topics"`
}

func LoadPredefinedTopics(file string) ([]PredefinedTopics, error) {
	topics := make([]PredefinedTopics, 0)

	yamlStr, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println("[Error] faild to read yaml file.")
		return nil, err
	}

	err = yaml.Unmarshal(yamlStr, &topics)
	if err != nil {
		log.Println("[Error] faild to unmarshall yaml file.")
		return nil, err
	}
	return topics, nil
}
