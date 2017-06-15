package ceratopoGon

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type PredefinedTopics map[string]map[string]uint16

func LoadPredefinedTopics(file string) (PredefinedTopics, error) {
	topics := PredefinedTopics{}

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
