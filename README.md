# gophercast

A simple, modular publish-subscribe (pub/sub) system written in Go with zero external dependencies.

## What is GopherCast?

GopherCast is a pub/sub messaging system where:

- **Publishers** send messages to topics.
- **Subscribers** receive messages from topics they're 'subscribed' to.
- **Brokers** manage the routing of messages from publishers to subscribers.

Publishers and subscribers are decoupled; they don't need to know about each other.

## Features

- Simple, clear code that is easily understood
- Zero external dependencies (Go standard library only)
- Modular design (each component can run independently)
- Thread-safe (safe for concurrent use)
- Non-blocking message delivery using goroutines
- In-process only (no network transport)

## Installation

### go get

```bash
go get github.com/gophercast/gophercast
```
## git

```bash
git clone github.com/gophercast/gophercast
```

## Quick Start

gophercast has a simple implementation you can easily build more features around:

```go
package main

import (
    "fmt"
    "github.com/gophercast/gophercast/internal/domain/broker"
    "github.com/gophercast/gophercast/internal/domain/message"
    "github.com/gophercast/gophercast/internal/domain/topic"
)

func main() {
    // Create broker
    b := broker.NewBroker()
    defer b.Close()
    
    // Create topic
    topic, _ := topic.New("events")
    
    // Subscribe
    sub := b.Subscribe(topic)
    defer b.Unsubscribe(sub.ID())
    
    // Listen for messages
    go func() {
        for msg := range sub.MessageChannel() {
            fmt.Println("Received:", msg.Data())
        }
    }()
    
    // Publish
    msg := message.NewMessage(topic, "Hello, World!")
    b.Publish(msg)
}
```
### Concepts:

**Topic**: A named channel (e.g., "user.created", "order.placed").

**Message**: Data being sent (includes topic, data, timestamp).

**Subscription**: Registration to receive messages from a topic.

**Broker**: Central hub that manages distribution.


```
┌─────────────────────────────────────────┐
│              Broker                     │
│  - Manages topics and subscriptions     │
│  - Routes messages to subscribers       │
└──────────┬────────────────┬─────────────┘
           │                │
           ▼                ▼
    ┌─────────────┐   ┌─────────────┐
    │ Publisher   │   │ Subscriber  │
    │ (sends msg) │   │ (recv msg)  │
    └─────────────┘   └─────────────┘
```

## Usage Examples

### Example 1: Basic Pub/Sub

```go
package main

import (
	"fmt"
	"time"

	"github.com/gophercast/gophercast/internal/domain/broker"
	"github.com/gophercast/gophercast/internal/domain/message"
	"github.com/gophercast/gophercast/internal/domain/topic"
)

func main() {
	fmt.Println("=== GopherCast System Example ===")

	// Step 1: Create broker
	fmt.Println("1. Creating broker...")
	b := broker.NewBroker()
	defer b.Close()

	// Step 2: Create topics
	fmt.Println("2. Creating topics...")
	usersTopic, _ := topic.New("users")
	ordersTopic, _ := topic.New("orders")

	// Step 3: Create subscribers
	fmt.Println("3. Creating subscribers...")

	// Subscriber 1: Listens to users topic
	sub1 := b.Subscribe(usersTopic)
	defer b.Unsubscribe(sub1.ID())

	go func() {
		fmt.Println("   [Subscriber 1] Listening to 'users' topic...")
		for msg := range sub1.MessageChannel() {
			fmt.Printf("   [Subscriber 1] Received: %s\n", msg.String())
		}
	}()

	// Subscriber 2: Also listens to users topic
	sub2 := b.Subscribe(usersTopic)
	defer b.Unsubscribe(sub2.ID())

	go func() {
		fmt.Println("   [Subscriber 2] Listening to 'users' topic...")
		for msg := range sub2.MessageChannel() {
			fmt.Printf("   [Subscriber 2] Received: %s\n", msg.String())
		}
	}()

	// Subscriber 3: Listens to orders topic
	sub3 := b.Subscribe(ordersTopic)
	defer b.Unsubscribe(sub3.ID())

	go func() {
		fmt.Println("   [Subscriber 3] Listening to 'orders' topic...")
		for msg := range sub3.MessageChannel() {
			fmt.Printf("   [Subscriber 3] Received: %s\n", msg.String())
		}
	}()

	// Give subscribers time to start
	time.Sleep(100 * time.Millisecond)

	// Step 4: Publish messages
	fmt.Println("\n4. Publishing messages...")

	// Publish to users topic
	fmt.Println("   Publishing to 'users' topic...")
	msg1 := message.NewMessage(usersTopic, "User Alice created")
	b.Publish(msg1)

	time.Sleep(100 * time.Millisecond)

	msg2 := message.NewMessage(usersTopic, "User Bob created")
	b.Publish(msg2)

	time.Sleep(100 * time.Millisecond)

	// Publish to orders topic
	fmt.Println("\n   Publishing to 'orders' topic...")
	msg3 := message.NewMessage(ordersTopic, "Order #123 placed")
	b.Publish(msg3)

	// Wait for messages to be received
	time.Sleep(500 * time.Millisecond)

	fmt.Println("\n=== Example Complete ===")
	fmt.Println("\nObservations:")
	fmt.Println("- Subscribers 1 and 2 both received messages from 'users' topic")
	fmt.Println("- Subscriber 3 only received messages from 'orders' topic")
	fmt.Println("- Each subscriber got their own copy of the messages")
}
```

### Example 2: Multiple Topics

```go
b := broker.NewBroker()

userTopic, _ := topic.New("users")
orderTopic, _ := topic.New("orders")

userSub := b.Subscribe(userTopic)
orderSub := b.Subscribe(orderTopic)

// Publish to different topics
b.Publish(message.NewMessage(userTopic, "User created"))
b.Publish(message.NewMessage(orderTopic, "Order placed"))
```

### Example 3: Multiple Subscribers

```go
b := broker.NewBroker()
topic, _ := topic.New("events")

// Multiple subscribers to same topic
sub1 := b.Subscribe(topic)
sub2 := b.Subscribe(topic)
sub3 := b.Subscribe(topic)

// All three receive the same message
b.Publish(message.NewMessage(topic, "Event occurred"))
```

## Running Examples

```bash
# Run the complete example
go run examples/basic/main.go

# Run standalone broker
go run cmd/broker/main.go

# Run publisher example
go run cmd/publisher/main.go

# Run subscriber example
go run cmd/subscriber/main.go
```
