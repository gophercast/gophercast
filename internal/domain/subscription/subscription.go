package subscription

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/gophercast/gophercast/internal/domain/message"
	"github.com/gophercast/gophercast/internal/domain/topic"
)

// Subscription represents a subscriber's registration to receive messages from a topic.
// Each subscription has its own channel for receiving messages.
type Subscription struct {
	id             string
	topic          topic.Topic
	messageChannel chan message.Message
	createdAt      time.Time
	closed         bool
	mu             sync.Mutex
}

// NewSubscription creates a new subscription for the given topic.
// The subscription includes a buffered channel for receiving messages.
func NewSubscription(t topic.Topic) *Subscription {
	return &Subscription{
		id:             generateSubscriptionID(),
		topic:          t,
		messageChannel: make(chan message.Message, 200), // Large buffer for concurrent publishing
		createdAt:      time.Now(),
	}
}

// ID returns the unique subscription identifier.
func (s *Subscription) ID() string {
	return s.id
}

// Topic returns the topic this subscription is for.
func (s *Subscription) Topic() topic.Topic {
	return s.topic
}

// CreatedAt returns when the subscription was created.
func (s *Subscription) CreatedAt() time.Time {
	return s.createdAt
}

// MessageChannel returns the channel for receiving messages.
// Subscribers should read from this channel to receive messages.
func (s *Subscription) MessageChannel() <-chan message.Message {
	return s.messageChannel
}

// SendMessage attempts to send a message to the subscriber.
// This is non-blocking; if the channel is full or the subscription is closed, the message is dropped.
func (s *Subscription) SendMessage(msg message.Message) {
	s.mu.Lock()

	if s.closed {
		s.mu.Unlock()
		return
	}

	// Try non-blocking send
	select {
	case s.messageChannel <- msg:
		// Message sent successfully
	default:
		// Channel full, drop message (best-effort delivery)
	}

	s.mu.Unlock()
}

// Close closes the message channel.
// After closing, no more messages can be sent to this subscription.
func (s *Subscription) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return
	}

	s.closed = true
	close(s.messageChannel)
}

// generateSubscriptionID creates a unique identifier for a subscription.
func generateSubscriptionID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
