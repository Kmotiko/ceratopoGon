package ceratopoGon

import (
	"strings"
	"sync"
)

var topicEntry *Topics = &Topics{
	topicRoot:  NewTopicNode(),
	topicIndex: make(map[string][]*MqttSnSession, 0)}

func GetTopicEntry() *Topics {
	return topicEntry
}

// Structure to management topics
type Topics struct {
	mutex      sync.RWMutex
	topicRoot  *TopicNode
	topicIndex map[string][]*MqttSnSession // fixed topic. map topicname to subscribers
}

// Topic Node of Topic Tree
type TopicNode struct {
	// mutex       sync.RWMutex
	children    map[string]*TopicNode
	subscribers []*MqttSnSession
}

// ctor
func NewTopicNode() *TopicNode {
	node := &TopicNode{}
	return node
}

// check whether topic is wildcarded or not
func IsWildCarded(topic []string) bool {
	for _, item := range topic {
		if item == "#" || item == "+" {
			return true
		}
	}
	return false
}

// search subscribers
func (t *Topics) GetSubscriber(topic string) []*MqttSnSession {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	// search from index
	subscribers := t.topicIndex[topic]
	if subscribers == nil {
		subscribers = make([]*MqttSnSession, 0)
	}

	return subscribers
}

// append subscribers
// @return amount of subscriber
func (t *Topics) AppendSubscriber(
	topic string,
	subscriber *MqttSnSession) (subscribers []*MqttSnSession) {
	// Lock
	t.mutex.Lock()
	defer t.mutex.Unlock()

	// split topic
	topics := strings.Split(topic, "/")

	// if is wildcarded
	if wildcarded := IsWildCarded(topics); wildcarded == true {
		node := t.topicRoot
		// visit node
		for i, topic := range topics {
			next := node.children[topic]

			// add node
			if next == nil {
				next = NewTopicNode()
				node.children[topic] = next
			}
			node = next

			if (len(topics) - 1) == i {
				// add subscriber
				node.subscribers = append(node.subscribers, subscriber)
				subscribers = node.subscribers
			}
		}

	} else {
		// add subscriber
		subscribers = t.topicIndex[topic]
		if subscribers == nil {
			subscribers = make([]*MqttSnSession, 0)
		}
		subscribers = append(subscribers, subscriber)
		t.topicIndex[topic] = subscribers
	}

	return
}

func (t *Topics) RemoveSubscriber(topic string, subscriber *MqttSnSession) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	// TODO: implement
}

// func add topic
// @return amount of subscriber, error if problem is occured
// func (t *TopicNode) AppendChild() (int, error) {
// }
