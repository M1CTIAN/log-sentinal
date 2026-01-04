package main

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"strings"
)

func main() {
	// 1. SETUP KAFKA READER (The Consumer)
	// We act as part of a "Consumer Group" named "log-sentinel-group".
	// This means if we start 10 copies of this program, they will split the work.
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"localhost:9092"},
		Topic:       "logs",
		GroupID:     "log-sentinel-group", // Important: Identifies who we are
		MinBytes:    10e3,                 // 10KB
		MaxBytes:    10e6,                 // 10MB
		StartOffset: kafka.FirstOffset,    // Read from the beginning of history
	})
	defer reader.Close()

	fmt.Println("Analyst Service Started... filtering for threats.")

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			break
		}

		message := string(m.Value)

		// INTELLIGENCE LOGIC
		if strings.Contains(message, "ALERT") {
			// Print in RED color (ANSI escape codes)
			fmt.Printf("\n🚨 THREAT DETECTED: %s\n", message)
		} else {
			// Print normal logs faintly (optional, or just ignore them)
			fmt.Printf("\nv checking: %s", message)
		}
	}
}
