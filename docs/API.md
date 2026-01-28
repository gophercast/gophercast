# API Documentation

## Broker

### `broker.NewBroker() *Broker`

Creates a new message broker.

**Returns**: A new Broker instance ready to manage topics and subscriptions.

**Example**:
```go
b := broker.NewBroker()
defer b.Close()
```

---

### `broker.Subscribe(topic topic.Topic) *subscription.Subscription`

Creates a new subscription for the given topic.

**Parameters**:
- `topic`: The Topic to subscribe to

**Returns**: A new Subscription that can receive messages from the topic

**Example**:
```go
topic, _ := topic.New("users")
sub := b.Subscribe(topic)
```

---

### `broker.Unsubscribe(subscriptionID string)`

Removes a subscription from the broker.

**Parameters**:
- `subscriptionID`: The unique identifier of the subscription to remove

**Example**:
```go
b.Unsubscribe(sub.ID())
```

---

### `broker.Publish(msg message.Message)`

Sends a message to all subscribers of the message's topic.

**Parameters**:
- `msg`: The Message to publish

**Behavior**:
- Messages are delivered asynchronously using goroutines
- If no subscribers exist for the topic, the message is dropped
- If a subscriber's channel is full, the message is dropped (best-effort delivery)

**Example**:
```go
msg := message.NewMessage(topic, "Hello, World!")
b.Publish(msg)
```

---

### `broker.Close()`

Closes all subscriptions and shuts down the broker.

**Behavior**:
- Closes all subscription channels
- Clears the subscriptions map
- After closing, the broker should not be used

**Example**:
```go
b.Close()
```

---

## Topic

### `topic.New(name string) (Topic, error)`

Creates a new topic with the given name.

**Parameters**:
- `name`: The topic name (must contain only letters, numbers, dots, and hyphens)

**Returns**:
- A new Topic if validation succeeds
- An error if the name is empty or contains invalid characters

**Example**:
```go
topic, err := topic.New("user.created")
if err != nil {
    // Handle error
}
```

**Valid Topic Names**:
- `"users"`
- `"user.created"`
- `"system-status"`

**Invalid Topic Names**:
- `""` (empty)
- `"user created"` (contains space)
- `"user/created"` (contains slash)

---

### `topic.String() string`

Returns the topic name.

**Returns**: The string representation of the topic

**Example**:
```go
name := topic.String() // "user.created"
```

---

### `topic.Equals(other Topic) bool`

Compares two topics for equality.

**Parameters**:
- `other`: Another topic to compare

**Returns**: `true` if both topics have the same name, `false` otherwise

**Example**:
```go
if topic1.Equals(topic2) {
    // Topics are equal
}
```

---

## Message

### `message.NewMessage(topic topic.Topic, data interface{}) Message`

Creates a new message for the given topic with the provided data.

**Parameters**:
- `topic`: The topic this message belongs to
- `data`: The message payload (can be any type)

**Returns**: A new Message with:
- A unique ID automatically generated
- The current timestamp

**Example**:
```go
msg := message.NewMessage(topic, map[string]string{"id": "123"})
```

---

### `message.ID() string`

Returns the unique message identifier.

**Returns**: The message's unique ID

**Example**:
```go
id := msg.ID()
```

---

### `message.Topic() topic.Topic`

Returns the topic this message belongs to.

**Returns**: The message's topic

**Example**:
```go
topic := msg.Topic()
```

---

### `message.Data() interface{}`

Returns the message payload.

**Returns**: The message data

**Example**:
```go
data := msg.Data()
```

---

### `message.PublishedAt() time.Time`

Returns when the message was created.

**Returns**: The message's timestamp

**Example**:
```go
timestamp := msg.PublishedAt()
```

---

### `message.String() string`

Returns a human-readable representation of the message.

**Returns**: String in format "Message[{ID}] on topic[{topic}] at {timestamp}"

**Example**:
```go
str := msg.String() // "Message[abc123] on topic[users] at 2025-01-27T10:00:00Z"
```

---

## Subscription

### `subscription.NewSubscription(topic topic.Topic) *Subscription`

Creates a new subscription for the given topic.

**Parameters**:
- `topic`: The topic to subscribe to

**Returns**: A new Subscription with:
- A unique ID
- A buffered channel for receiving messages (buffer size: 200)

**Example**:
```go
sub := subscription.NewSubscription(topic)
```

---

### `subscription.ID() string`

Returns the unique subscription identifier.

**Returns**: The subscription's unique ID

**Example**:
```go
id := sub.ID()
```

---

### `subscription.Topic() topic.Topic`

Returns the topic this subscription is for.

**Returns**: The subscription's topic

**Example**:
```go
topic := sub.Topic()
```

---

### `subscription.MessageChannel() <-chan message.Message`

Returns the channel for receiving messages.

**Returns**: A receive-only channel that delivers messages to the subscriber

**Example**:
```go
for msg := range sub.MessageChannel() {
    // Process message
}
```

---

### `subscription.SendMessage(msg message.Message)`

Attempts to send a message to the subscriber.

**Parameters**:
- `msg`: The message to send

**Behavior**:
- Non-blocking send (uses select with default)
- If the channel is full, the message is dropped
- If the subscription is closed, the message is dropped

**Example**:
```go
sub.SendMessage(msg)
```

---

### `subscription.Close()`

Closes the message channel.

**Behavior**:
- Closes the underlying channel
- After closing, no more messages can be sent
- Safe to call multiple times (idempotent)

**Example**:
```go
sub.Close()
```

---

### `subscription.CreatedAt() time.Time`

Returns when the subscription was created.

**Returns**: The subscription's creation timestamp

**Example**:
```go
timestamp := sub.CreatedAt()
```
