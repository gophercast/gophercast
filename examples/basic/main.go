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
