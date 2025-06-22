package queue

import "log/slog"

type Event interface {}

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
}

func NewInMemoryQueue() *InMemoryQueue {
	return &InMemoryQueue{
		subscribers: make(map[string][]chan Event),
	}
}

func (q *InMemoryQueue) Subscribe(topic string) <-chan Event {
	ch := make(chan Event, 10) // Buffered channel to avoid blocking
	q.subscribers[topic] = append(q.subscribers[topic], ch)
	return ch
}

func (q *InMemoryQueue) Publish(topic string, message Event) {
	if subs, ok := q.subscribers[topic]; ok {
		for _, sub := range subs {
			select {
			case sub <- message:
			default:
				// If the channel is full, we skip sending to avoid blocking
				slog.Warn("Queue publish skipped", "topic", topic, "message", message)
			}
		}
	}
}

func (q *InMemoryQueue) Unsubscribe(topic string, ch <-chan Event) {
	if subs, ok := q.subscribers[topic]; ok {
		for i, sub := range subs {
			if sub == ch {
				q.subscribers[topic] = append(subs[:i], subs[i+1:]...)
				close(sub) // Close the channel to avoid memory leaks
				break
			}
		}
	}
}
