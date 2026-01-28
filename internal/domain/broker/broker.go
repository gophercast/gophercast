package broker

import (
	"sync"

	"github.com/gophercast/gophercast/internal/domain/message"
	"github.com/gophercast/gophercast/internal/domain/subscription"
	"github.com/gophercast/gophercast/internal/domain/topic"
)

// Broker is the central hub that manages topics and routes messages to subscribers.
// It is safe for concurrent use by multiple goroutines.
type Broker struct {
	subscriptions map[string][]*subscription.Subscription // topic name -> subscriptions
	mutex         sync.RWMutex
}

// NewBroker creates a new message broker.
func NewBroker() *Broker {
	return &Broker{
		subscriptions: make(map[string][]*subscription.Subscription),
	}
}

// Subscribe creates a new subscription for the given topic.
// Returns the subscription which includes a channel for receiving messages.
func (b *Broker) Subscribe(t topic.Topic) *subscription.Subscription {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	sub := subscription.NewSubscription(t)
	topicName := t.String()

	b.subscriptions[topicName] = append(b.subscriptions[topicName], sub)

	return sub
}

// Unsubscribe removes a subscription from the broker.
// The subscription will no longer receive messages.
func (b *Broker) Unsubscribe(subscriptionID string) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Find and remove the subscription from all topics
	for topicName, subs := range b.subscriptions {
		for i, sub := range subs {
			if sub.ID() == subscriptionID {
				// Close the subscription
				sub.Close()

				// Remove from slice
				b.subscriptions[topicName] = append(subs[:i], subs[i+1:]...)

				// If no more subscriptions for this topic, remove the topic
				if len(b.subscriptions[topicName]) == 0 {
					delete(b.subscriptions, topicName)
				}

				return
			}
		}
	}
}

// Publish sends a message to all subscribers of the message's topic.
// Distribution is done concurrently using goroutines to avoid blocking.
// Messages sent to topics with no subscribers are dropped.
func (b *Broker) Publish(msg message.Message) {
	b.mutex.RLock()
	subs := b.subscriptions[msg.Topic().String()]
	b.mutex.RUnlock()

	// If no subscribers, message is dropped
	if len(subs) == 0 {
		return
	}

	// Send to all subscribers concurrently
	for _, sub := range subs {
		go sub.SendMessage(msg)
	}
}

// Close closes all subscriptions and shuts down the broker.
// After closing, the broker should not be used.
func (b *Broker) Close() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Close all subscriptions
	for _, subs := range b.subscriptions {
		for _, sub := range subs {
			sub.Close()
		}
	}

	// Clear the subscriptions map
	b.subscriptions = make(map[string][]*subscription.Subscription)
}
