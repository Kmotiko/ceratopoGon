package ceratopogon

import (
	"sync"
)

type TopicMap struct {
	mutex   sync.RWMutex
	topics  map[uint16]string
	topicId uint16
}

func NewTopicMap() *TopicMap {
	t := &TopicMap{
		topics:  make(map[uint16]string),
		topicId: 0}
	return t
}

func (t *TopicMap) NextTopicId() uint16 {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.topicId++
	if t.topicId == 0xffff {
		// err
	}

	return t.topicId
}

func (t *TopicMap) StoreTopic(topicName string) uint16 {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	topicId, ok := t.LoadTopicId(topicName)
	if ok {
		return topicId
	}

	topicId = t.NextTopicId()

	t.topics[topicId] = topicName
	return topicId
}

func (t *TopicMap) LoadTopic(topicId uint16) (string, bool) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	topicName, ok := t.topics[topicId]
	return topicName, ok
}

func (t *TopicMap) LoadTopicId(topicName string) (uint16, bool) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	for k, v := range t.topics {
		if v == topicName {
			return k, true
		}
	}
	return 0, false
}
