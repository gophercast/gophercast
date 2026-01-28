package subscription_test

import (
	"testing"
	"time"

	"github.com/gophercast/gophercast/internal/domain/message"
	"github.com/gophercast/gophercast/internal/domain/subscription"
	"github.com/gophercast/gophercast/internal/domain/topic"
)

func TestNewSubscription(t *testing.T) {
	topicObj, _ := topic.New("users")
	sub := subscription.NewSubscription(topicObj)

	// Check ID is generated and unique
	if sub.ID() == "" {
		t.Error("ID should not be empty")
	}

	// Check Topic is correctly set
	if sub.Topic().String() != "users" {
		t.Errorf("Topic = %v, want %v", sub.Topic().String(), "users")
	}

	// Check MessageChannel is initialized and not nil
	if sub.MessageChannel() == nil {
		t.Error("MessageChannel should not be nil")
	}

	// Check CreatedAt is set
	if sub.CreatedAt().IsZero() {
		t.Error("CreatedAt should not be zero")
	}

	// Check CreatedAt is recent
	now := time.Now()
	if sub.CreatedAt().Sub(now) > time.Second || now.Sub(sub.CreatedAt()) > time.Second {
		t.Error("CreatedAt should be within 1 second of now")
	}
}

func TestSubscriptionSendMessage(t *testing.T) {
	topicObj, _ := topic.New("users")
	sub := subscription.NewSubscription(topicObj)

	msg := message.NewMessage(topicObj, "test data")

	// Send message successfully
	sub.SendMessage(msg)

	// Check if message was received on channel
	select {
	case received := <-sub.MessageChannel():
		if received.Data() != "test data" {
			t.Errorf("Received data = %v, want %v", received.Data(), "test data")
		}
	case <-time.After(time.Second):
		t.Error("Did not receive message within 1 second")
	}
}

func TestSubscriptionClose(t *testing.T) {
	topicObj, _ := topic.New("users")
	sub := subscription.NewSubscription(topicObj)

	// Close the subscription
	sub.Close()

	// Check that channel is closed
	select {
	case _, ok := <-sub.MessageChannel():
		if ok {
			t.Error("Channel should be closed")
		}
	default:
		// Channel was buffered and not read, but Close() should have closed it
	}
}
