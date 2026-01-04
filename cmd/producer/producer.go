package main

import (
	"context"
	"fmt"
	"os"
	"time"
	"math/rand"
	"github.com/nxadm/tail"
	"github.com/segmentio/kafka-go"
)

func main() {
	// 1. SETUP KAFKA CONNECTION
	// We connect to "localhost:9092" because that is the port we opened in docker-compose.yml
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "logs",
		Balancer: &kafka.LeastBytes{}, // Distributes messages evenly
	})
	defer writer.Close()

	// 2. FILE SETUP
	// We create a dummy file to watch so you can see it working immediately
	fileName := "app.log"
	createDummyFile(fileName)

	// 3. START THE SPY (TAIL)
	config := tail.Config{
		Follow: true, // Keep looking for new lines
		ReOpen: true, // If file gets deleted, reopen it
	}
	t, err := tail.TailFile(fileName, config)
	if err != nil {
		panic(err)
	}

	fmt.Println("🕵️  Spy Agent Started... Forwarding logs to Kafka")

	// 4. GENERATE FAKE LOGS (Background Job)
	// This simulates a real website writing logs
	go generateFakeLogs(fileName)

	// 5. MAIN LOOP: Read File -> Send to Kafka
	for line := range t.Lines {
		// Print to screen (so you know it's happening)
		fmt.Println("Sending:", line.Text)

		// Send to Kafka
		err := writer.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte(fmt.Sprintf("Key-%d", time.Now().Unix())),
				Value: []byte(line.Text),
			},
		)

		if err != nil {
			fmt.Println("❌ Failed to write to Kafka:", err)
		}
	}
}

// --- HELPER FUNCTIONS ---

func createDummyFile(fileName string) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		file, _ := os.Create(fileName)
		file.Close()
	}
}

func generateFakeLogs(fileName string) {
	f, _ := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
	defer f.Close()

	for {
		time.Sleep(1 * time.Second)

		// 10% chance to generate a hacking log
		var logMessage string
		if rand.Intn(2) == 0 {
			logMessage = fmt.Sprintf("ALERT: Failed login attempt from IP 192.168.1.%d at %s\n", rand.Intn(255), time.Now().Format("15:04:05"))
		} else {
			logMessage = fmt.Sprintf("INFO: User action detected at %s\n", time.Now().Format("15:04:05"))
		}
		
		f.WriteString(logMessage)
	}
}