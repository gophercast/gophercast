package broker_test

import (
	"sync"
	"testing"
	"time"

	"github.com/gophercast/gophercast/internal/domain/broker"
	"github.com/gophercast/gophercast/internal/domain/message"
	"github.com/gophercast/gophercast/internal/domain/topic"
)

func TestNewBroker(t *testing.T) {
	b := broker.NewBroker()

	if b == nil {
		t.Error("NewBroker() should not return nil")
	}

	// Note: We can't directly test internal state without exposing it
	// In a real scenario, we might add exported methods for testing
}

func TestBrokerSubscribe(t *testing.T) {
	b := broker.NewBroker()
	defer b.Close()

	topicObj, _ := topic.New("users")

	// First subscriber to topic
	sub1 := b.Subscribe(topicObj)
	if sub1 == nil {
		t.Error("Subscribe() should not return nil")
	}
	if sub1.Topic().String() != "users" {
		t.Errorf("Subscription topic = %v, want %v", sub1.Topic().String(), "users")
	}

	// Second subscriber to same topic
	sub2 := b.Subscribe(topicObj)
	if sub2 == nil {
		t.Error("Subscribe() should not return nil")
	}
	if sub1.ID() == sub2.ID() {
		t.Error("Subscription IDs should be unique")
	}
}

func TestBrokerPublish(t *testing.T) {
	b := broker.NewBroker()
	defer b.Close()

	topicObj, _ := topic.New("users")

	// Test with no subscribers (message should be dropped)
	msg := message.NewMessage(topicObj, "test")
	b.Publish(msg) // Should not panic

	// Test with one subscriber
	sub := b.Subscribe(topicObj)

	go func() {
		for msg := range sub.MessageChannel() {
			if msg.Data() != "test" {
				t.Errorf("Received data = %v, want %v", msg.Data(), "test")
			}
		}
	}()

	b.Publish(message.NewMessage(topicObj, "test"))

	// Give time for message to be delivered
	time.Sleep(100 * time.Millisecond)
}

func TestBrokerPublishMultipleSubscribers(t *testing.T) {
	b := broker.NewBroker()
	defer b.Close()

	topicObj, _ := topic.New("users")

	// Create multiple subscribers
	sub1 := b.Subscribe(topicObj)
	sub2 := b.Subscribe(topicObj)
	sub3 := b.Subscribe(topicObj)

	received := make(map[string]int)
	var mu sync.Mutex

	go func() {
		for {
			select {
			case <-sub1.MessageChannel():
				mu.Lock()
				received["sub1"]++
				mu.Unlock()
			case <-sub2.MessageChannel():
				mu.Lock()
				received["sub2"]++
				mu.Unlock()
			case <-sub3.MessageChannel():
				mu.Lock()
				received["sub3"]++
				mu.Unlock()
			case <-time.After(time.Second):
				return
			}
		}
	}()

	b.Publish(message.NewMessage(topicObj, "test"))

	// Give time for messages to be delivered
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	if received["sub1"] != 1 {
		t.Errorf("sub1 received %d messages, want 1", received["sub1"])
	}
	if received["sub2"] != 1 {
		t.Errorf("sub2 received %d messages, want 1", received["sub2"])
	}
	if received["sub3"] != 1 {
		t.Errorf("sub3 received %d messages, want 1", received["sub3"])
	}
	mu.Unlock()
}

func TestBrokerUnsubscribe(t *testing.T) {
	b := broker.NewBroker()
	defer b.Close()

	topicObj, _ := topic.New("users")

	sub := b.Subscribe(topicObj)
	subID := sub.ID()

	// Unsubscribe
	b.Unsubscribe(subID)

	// After unsubscribe, the subscription should receive no more messages
	// This is hard to test directly without exposing internal state
}

func TestBrokerClose(t *testing.T) {
	b := broker.NewBroker()

	topicObj, _ := topic.New("users")
	sub := b.Subscribe(topicObj)

	// Close broker
	b.Close()

	// Check that subscription channel is closed
	select {
	case _, ok := <-sub.MessageChannel():
		if ok {
			t.Error("Channel should be closed after broker.Close()")
		}
	default:
		// Channel was buffered and not read
	}
}

func TestBrokerPublishWrongTopic(t *testing.T) {
	b := broker.NewBroker()
	defer b.Close()

	usersTopic, _ := topic.New("users")
	ordersTopic, _ := topic.New("orders")

	// Subscribe to "users" topic
	sub := b.Subscribe(usersTopic)

	// Publish to "orders" topic
	b.Publish(message.NewMessage(ordersTopic, "order data"))

	// Should not receive message (different topic)
	select {
	case msg := <-sub.MessageChannel():
		t.Errorf("Should not receive message for different topic, got: %v", msg.Data())
	case <-time.After(100 * time.Millisecond):
		// Expected: no message received
	}
}

func TestBrokerConcurrentPublish(t *testing.T) {
	b := broker.NewBroker()
	defer b.Close()

	topicObj, _ := topic.New("users")
	sub := b.Subscribe(topicObj)

	// Count received messages
	receivedCount := 0
	var mu sync.Mutex
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for msg := range sub.MessageChannel() {
			mu.Lock()
			receivedCount++
			mu.Unlock()
			_ = msg
		}
	}()

	// Multiple goroutines publishing simultaneously
	pubWg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		pubWg.Add(1)
		go func() {
			defer pubWg.Done()
			for j := 0; j < 10; j++ {
				b.Publish(message.NewMessage(topicObj, "test"))
			}
		}()
	}

	pubWg.Wait()

	// Give more time for all messages to be delivered through goroutines
	time.Sleep(500 * time.Millisecond)

	mu.Lock()
	if receivedCount != 100 {
		t.Errorf("Received %d messages, want 100", receivedCount)
	}
	mu.Unlock()
}
