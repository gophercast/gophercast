package main

import (
	"fmt"
	"time"

	"github.com/gophercast/gophercast/internal/domain/broker"
	"github.com/gophercast/gophercast/internal/domain/message"
	"github.com/gophercast/gophercast/internal/domain/topic"
)

func main() {
	fmt.Println("Starting Publisher Example...")

	// Create broker (in real use, this would be shared)
	b := broker.NewBroker()
	defer b.Close()

	// Create topic
	usersTopic, err := topic.New("users.created")
	if err != nil {
		fmt.Printf("Error creating topic: %v\n", err)
		return
	}

	// Publish messages
	fmt.Println("\nPublishing messages...")

	for i := 1; i <= 5; i++ {
		// Create message data
		data := map[string]interface{}{
			"user_id": fmt.Sprintf("user-%d", i),
			"email":   fmt.Sprintf("user%d@example.com", i),
		}

		// Create and publish message
		msg := message.NewMessage(usersTopic, data)
		b.Publish(msg)

		fmt.Printf("Published: %s with data: %v\n", msg.String(), data)

		// Wait a bit between messages
		time.Sleep(500 * time.Millisecond)
	}

	fmt.Println("\nAll messages published!")

	// Keep running for a bit to let subscribers receive
	time.Sleep(2 * time.Second)
}
