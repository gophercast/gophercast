package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gophercast/gophercast/internal/domain/broker"
	"github.com/gophercast/gophercast/internal/domain/topic"
)

func main() {
	fmt.Println("Starting Subscriber Example...")

	// Create broker (in real use, this would be shared)
	b := broker.NewBroker()
	defer b.Close()

	// Create topic to subscribe to
	usersTopic, err := topic.New("users.created")
	if err != nil {
		fmt.Printf("Error creating topic: %v\n", err)
		return
	}

	// Subscribe to topic
	sub := b.Subscribe(usersTopic)
	defer b.Unsubscribe(sub.ID())

	fmt.Printf("Subscribed to topic: %s\n", usersTopic.String())
	fmt.Println("Waiting for messages... Press Ctrl+C to stop.")

	// Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Listen for messages
	go func() {
		for msg := range sub.MessageChannel() {
			fmt.Printf("\nReceived: %s\n", msg.String())
			fmt.Printf("Data: %v\n", msg.Data())
		}
	}()

	// Wait for shutdown signal
	<-quit
	fmt.Println("\nShutting down subscriber...")
}
