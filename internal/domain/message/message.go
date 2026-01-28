package message

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gophercast/gophercast/internal/domain/topic"
)

// Message represents data being published to a topic.
// Messages are immutable once created.
type Message struct {
	id          string
	topic       topic.Topic
	data        interface{}
	publishedAt time.Time
}

// NewMessage creates a new message for the given topic with the provided data.
// A unique ID and timestamp are automatically assigned.
func NewMessage(t topic.Topic, data interface{}) Message {
	return Message{
		id:          generateMessageID(),
		topic:       t,
		data:        data,
		publishedAt: time.Now(),
	}
}

// ID returns the unique message identifier.
func (m Message) ID() string {
	return m.id
}

// Topic returns the topic this message belongs to.
func (m Message) Topic() topic.Topic {
	return m.topic
}

// Data returns the message payload.
func (m Message) Data() interface{} {
	return m.data
}

// PublishedAt returns when the message was created.
func (m Message) PublishedAt() time.Time {
	return m.publishedAt
}

// String returns a human-readable representation of the message.
func (m Message) String() string {
	return fmt.Sprintf("Message[%s] on topic[%s] at %s",
		m.id, m.topic.String(), m.publishedAt.Format(time.RFC3339))
}

// generateMessageID creates a unique identifier for a message.
func generateMessageID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
