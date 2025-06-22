package queue

import (
	"log/slog"
	"sync"
)

type Event interface{}

type BaseEvent struct {
	EventData interface{} // The actual event data
}

func (e BaseEvent) Event() interface{} {
	return e.EventData
}

func NewEvent(eventData interface{}) Event {
	return BaseEvent{EventData: eventData}
}

type Queue interface {
	Subscribe(topic string) <-chan Event
	Unsubscribe(topic string, ch <-chan Event)
	Publish(topic string, message Event)
}

type InMemoryQueue struct {
	subscribers map[string][]chan Event
	mu          sync.RWMutex
}

func NewInMemoryQueue() *InMemoryQueue {
	return &InMemoryQueue{
		subscribers: make(map[string][]chan Event),
	}
}

func (q *InMemoryQueue) Subscribe(topic string) <-chan Event {
	ch := make(chan Event, 50) // Buffered channel to avoid blocking

	q.mu.Lock()
	q.subscribers[topic] = append(q.subscribers[topic], ch)
	q.mu.Unlock()

	return ch
}

func (q *InMemoryQueue) Publish(topic string, message Event) {
	q.mu.RLock()
	subs, ok := q.subscribers[topic]
	if !ok {
		q.mu.RUnlock()
		return
	}

	// 创建 subs 的副本，避免在读锁期间长时间持有
	subsCopy := make([]chan Event, len(subs))
	copy(subsCopy, subs)
	q.mu.RUnlock()

	for _, sub := range subsCopy {
		select {
		case sub <- message:
		default:
			// If the channel is full, we skip sending to avoid blocking
			slog.Warn("Queue publish skipped", "topic", topic, "message", message)
		}
	}
}

func (q *InMemoryQueue) Unsubscribe(topic string, ch <-chan Event) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if subs, ok := q.subscribers[topic]; ok {
		for i, sub := range subs {
			if sub == ch {
				q.subscribers[topic] = append(subs[:i], subs[i+1:]...)
				close(sub) // Close the channel to avoid memory leaks
				break
			}
		}

		// 如果该 topic 没有订阅者了，删除该 topic
		if len(q.subscribers[topic]) == 0 {
			delete(q.subscribers, topic)
		}
	}
}
