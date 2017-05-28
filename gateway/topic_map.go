package ceratopogon

import (
	"log"
	"sync"
)

type TopicMap struct {
	mutex   sync.RWMutex
	topics  map[uint16]string
	topicId *TopicId
}

func NewTopicMap() *TopicMap {
	t := &TopicMap{
		mutex:   sync.RWMutex{},
		topics:  make(map[uint16]string),
		topicId: &TopicId{}}
	return t
}

func (t *TopicMap) StoreTopicWithId(topicName string, id uint16) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	ok := t.topicId.EnsureId(id)
	if !ok {
		log.Println("[Error] failed to ensure topic id.")
		return false
	}

	t.topics[id] = topicName
	return true
}

func (t *TopicMap) StoreTopic(topicName string) uint16 {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for k, v := range t.topics {
		if v == topicName {
			return k
		}
	}

	// get next id
	id := t.topicId.NextTopicId()

	t.topics[id] = topicName
	return id
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

type TopicId struct {
	mutex   sync.RWMutex
	topicId [0xffff]bool
}

func (t *TopicId) EnsureId(id uint16) bool {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if !t.topicId[id] {
		t.topicId[id] = true
		return true
	}

	return false
}

func (t *TopicId) NextTopicId() uint16 {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for i, v := range t.topicId {
		if v == false {
			t.topicId[i] = true
			return uint16(i)
		}
	}

	// TODO: error handling

	return 0
}

func (t *TopicId) FreeTopicId(i uint16) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.topicId[i] = false

	return
}
