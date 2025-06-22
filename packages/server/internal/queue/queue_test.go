package queue

import (
	"sync"
	"testing"
	"time"
)

func TestInMemoryQueue_BasicOperations(t *testing.T) {
	q := NewInMemoryQueue()

	t.Run("subscribe_and_publish", func(t *testing.T) {
		topic := "test-topic"
		ch := q.Subscribe(topic)

		testEvent := "test-message"
		q.Publish(topic, testEvent)

		select {
		case received := <-ch:
			if received != testEvent {
				t.Errorf("Expected %v, got %v", testEvent, received)
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("Expected to receive message within timeout")
		}
	})

	t.Run("publish_to_nonexistent_topic", func(t *testing.T) {
		q.Publish("nonexistent-topic", "test-message")
		// Should not panic or cause issues
	})

	t.Run("unsubscribe", func(t *testing.T) {
		topic := "unsubscribe-test"
		ch := q.Subscribe(topic)

		q.Unsubscribe(topic, ch)

		q.Publish(topic, "test-message")

		// 验证 channel 已关闭
		select {
		case msg, ok := <-ch:
			if ok {
				t.Errorf("Should not receive message after unsubscribe, got: %v", msg)
			}
			// Channel 关闭是正确的行为
		case <-time.After(50 * time.Millisecond):
			t.Error("Expected channel to be closed after unsubscribe")
		}
	})
}

func TestInMemoryQueue_MultipleSubscribers(t *testing.T) {
	q := NewInMemoryQueue()
	topic := "multi-subscriber-test"

	const numSubscribers = 3
	channels := make([]<-chan Event, numSubscribers)

	for i := 0; i < numSubscribers; i++ {
		channels[i] = q.Subscribe(topic)
	}

	testEvent := "broadcast-message"
	q.Publish(topic, testEvent)

	for i, ch := range channels {
		select {
		case received := <-ch:
			if received != testEvent {
				t.Errorf("Subscriber %d: expected %v, got %v", i, testEvent, received)
			}
		case <-time.After(100 * time.Millisecond):
			t.Errorf("Subscriber %d: expected to receive message within timeout", i)
		}
	}
}

func TestInMemoryQueue_ConcurrentAccess(t *testing.T) {
	q := NewInMemoryQueue()

	const numGoroutines = 10
	const messagesPerGoroutine = 5

	var wg sync.WaitGroup

	// 并发订阅
	t.Run("concurrent_subscribe", func(t *testing.T) {
		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()

				for j := 0; j < messagesPerGoroutine; j++ {
					topic := "concurrent-topic"
					ch := q.Subscribe(topic)

					// 立即取消订阅，测试并发修改
					q.Unsubscribe(topic, ch)
				}
			}(i)
		}

		wg.Wait()
	})

	// 并发发布
	t.Run("concurrent_publish", func(t *testing.T) {
		topic := "publish-test"
		ch := q.Subscribe(topic)

		// 在单独的 goroutine 中消费消息，避免缓冲区满
		receivedCount := 0
		var receiveMu sync.Mutex
		done := make(chan struct{})

		go func() {
			for {
				select {
				case <-ch:
					receiveMu.Lock()
					receivedCount++
					receiveMu.Unlock()
				case <-done:
					return
				}
			}
		}()

		wg.Add(numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer wg.Done()

				for j := 0; j < messagesPerGoroutine; j++ {
					q.Publish(topic, id*messagesPerGoroutine+j)
				}
			}(i)
		}

		wg.Wait()

		// 等待一段时间让消息被消费
		time.Sleep(100 * time.Millisecond)
		close(done)

		receiveMu.Lock()
		finalCount := receivedCount
		receiveMu.Unlock()

		// 由于缓冲区限制，我们至少应该收到一些消息
		if finalCount == 0 {
			t.Error("Expected to receive at least some messages")
		}

		q.Unsubscribe(topic, ch)
	})
}

func TestInMemoryQueue_BufferedChannels(t *testing.T) {
	q := NewInMemoryQueue()
	topic := "buffer-test"
	ch := q.Subscribe(topic)

	// 发送超过缓冲区大小的消息
	const messageCount = 15 // 大于默认缓冲区大小 10

	for i := 0; i < messageCount; i++ {
		q.Publish(topic, i)
	}

	// 应该能接收到至少缓冲区大小的消息
	receivedCount := 0
	timeout := time.After(100 * time.Millisecond)

receiveLoop:
	for {
		select {
		case <-ch:
			receivedCount++
		case <-timeout:
			break receiveLoop
		}
	}

	if receivedCount < 10 {
		t.Errorf("Expected at least 10 messages in buffer, got %d", receivedCount)
	}

	q.Unsubscribe(topic, ch)
}

func TestInMemoryQueue_BaseEvent(t *testing.T) {
	testData := "test-string-data"

	event := NewEvent(testData)
	baseEvent, ok := event.(BaseEvent)
	if !ok {
		t.Fatal("Expected BaseEvent type")
	}

	if baseEvent.Event() != testData {
		t.Errorf("Expected %v, got %v", testData, baseEvent.Event())
	}
}

func BenchmarkInMemoryQueue_Subscribe(b *testing.B) {
	q := NewInMemoryQueue()
	topic := "bench-subscribe"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Subscribe(topic)
	}
}

func BenchmarkInMemoryQueue_Publish(b *testing.B) {
	q := NewInMemoryQueue()
	topic := "bench-publish"

	// 预先创建一些订阅者
	for i := 0; i < 5; i++ {
		q.Subscribe(topic)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Publish(topic, i)
	}
}

func BenchmarkInMemoryQueue_ConcurrentOperations(b *testing.B) {
	q := NewInMemoryQueue()
	topic := "bench-concurrent"

	// Pre-create some subscribers to avoid race conditions
	const numSubscribers = 5
	channels := make([]<-chan Event, numSubscribers)
	for i := 0; i < numSubscribers; i++ {
		channels[i] = q.Subscribe(topic)
	}

	// Consume messages in background to prevent buffer overflow
	done := make(chan struct{})
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				for _, ch := range channels {
					select {
					case <-ch:
					default:
					}
				}
				time.Sleep(time.Microsecond)
			}
		}
	}()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			q.Publish(topic, "test")
		}
	})

	close(done)

	// Clean up
	for _, ch := range channels {
		q.Unsubscribe(topic, ch)
	}
}
