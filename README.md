# Log Sentinel

Log Sentinel is a distributed, event-driven log monitoring system built with Go and Apache Kafka.

It demonstrates a decoupled microservices architecture designed to ingest, transport, and analyze log data in real-time. The system simulates a production environment where a "Spy" service watches application logs and forwards them to a Kafka cluster, while an "Analyst" service consumes the stream to detect potential security threats.

## Architecture

The project follows a Producer-Consumer pattern utilizing Kafka as the central message broker:

1.  **Log Generator (Simulation):** A background routine that generates realistic application logs, including random security alerts (e.g., failed login attempts).
2.  **Producer Service (The Spy):** A Go service that "tails" the log file in real-time. It captures new lines and publishes them to the Kafka `logs` topic.
3.  **Kafka Cluster:** Running in Docker using KRaft mode (no Zookeeper dependency) for lightweight infrastructure.
4.  **Consumer Service (The Analyst):** A Go service that subscribes to the `logs` topic. It processes the stream and triggers visual alerts when security threats are detected.

## Tech Stack

* **Language:** Go (Golang) 1.21+
* **Message Broker:** Apache Kafka (Confluent Image, KRaft Mode)
* **Infrastructure:** Docker & Docker Compose
* **Key Libraries:**
    * `segmentio/kafka-go`: Pure Go Kafka client.
    * `nxadm/tail`: Real-time file watching.

## Project Structure

```text
log-sentinel/
├── cmd/
│   ├── producer/
│   │   └── main.go    # The Log Watcher & Forwarder
│   └── consumer/
│       └── main.go    # The Log Analyzer & Alert System
├── docker-compose.yml     # Kafka Infrastructure Configuration
├── go.mod                 # Go Dependencies
└── README.md
