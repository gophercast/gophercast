package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gophercast/gophercast/internal/domain/broker"
)

func main() {
	fmt.Println("Starting GopherCast Broker...")

	// Create broker
	b := broker.NewBroker()
	defer b.Close()

	fmt.Println("GopherCast Broker is running. Press Ctrl+C to stop.")
	fmt.Println("Note: This is a demonstration. Real usage involves importing the broker in your code.")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down broker...")
}
