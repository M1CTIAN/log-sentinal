package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq" // The Postgres Driver
	"github.com/segmentio/kafka-go"
)

func main() {
	// 1. CONNECT TO DATABASE
	// Connection string: user=admin password=secretpassword dbname=logsentinel sslmode=disable
	connStr := "postgres://admin:secretpassword@localhost:5432/logsentinel?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test the connection
	if err := db.Ping(); err != nil {
		log.Fatal("❌ Could not connect to database:", err)
	}
	fmt.Println("✅ Connected to the Vault (PostgreSQL)")

	// 2. CREATE TABLE (Migration)
	// We create the table automatically so you don't have to use SQL manually
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS alerts (
		id SERIAL PRIMARY KEY,
		message TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(createTableQuery); err != nil {
		log.Fatal("Failed to create table:", err)
	}

	// 3. START KAFKA CONSUMER
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:9092"},
		Topic:     "logs",
		GroupID:   "log-sentinel-group",
		StartOffset: kafka.FirstOffset,
	})
	defer reader.Close()

	fmt.Println("🦅 Analyst Service Started... Filtering for threats.")

	// 4. PROCESS LOOP
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			break
		}

		message := string(m.Value)

		// INTELLIGENCE LOGIC
		if strings.Contains(message, "ALERT") {
			fmt.Printf("\n🚨 THREAT DETECTED: %s\n", message)
			
			// SAVE TO DB
			_, err := db.Exec("INSERT INTO alerts (message) VALUES ($1)", message)
			if err != nil {
				fmt.Println("❌ Failed to save to DB:", err)
			} else {
				fmt.Println("💾 Saved to Vault.")
			}
		} else {
			// Ignore normal logs
			fmt.Printf("\nNormal log received: %s\n", message)
		}
	}
}