package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/segmentio/kafka-go"
)

type Activity struct {
	UserID string `json:"user_id"`
	Action string `json:"action"`
}

var (
	kafkaBroker   = os.Getenv("KAFKA_BROKER")
	kafkaTopic    = os.Getenv("KAFKA_TOPIC")
	clickhouseAddr = os.Getenv("CLICKHOUSE_ADDR")
	db            clickhouse.Conn
)

func main() {
	var err error
	// Initialize ClickHouse connection
	db, err = connectToClickHouse()
	if err != nil {
		log.Fatalf("Failed to connect to ClickHouse: %v", err)
	}
	defer db.Close()
	log.Println("Worker successfully connected to ClickHouse")

	// Create table if it doesn't exist
	createTable()

	// Initialize Kafka Reader
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		Topic:   kafkaTopic,
		GroupID: "activity-workers",
	})
	defer kafkaReader.Close()
	log.Println("Kafka reader connected and waiting for messages...")

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	
	for {
		select {
		case <-sigchan:
			log.Println("Termination signal received, shutting down.")
			return
		default:
			msg, err := kafkaReader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("could not read message: %v", err)
				continue
			}

			var activity Activity
			if err := json.Unmarshal(msg.Value, &activity); err != nil {
				log.Printf("could not unmarshal message: %v", err)
				continue
			}

			log.Printf("Received message: UserID=%s, Action=%s", activity.UserID, activity.Action)
			saveActivity(activity)
		}
	}
}

func connectToClickHouse() (clickhouse.Conn, error) {
    // ... (This function is the same as in the api/main.go file)
}

func createTable() {
	err := db.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS activities (
			user_id String,
			action String,
			created_at DateTime DEFAULT now()
		) ENGINE = MergeTree()
		ORDER BY created_at;
	`)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	log.Println("Table 'activities' is ready.")
}

func saveActivity(activity Activity) {
	err := db.Exec(context.Background(), "INSERT INTO activities (user_id, action) VALUES (?, ?)",
		activity.UserID,
		activity.Action,
	)
	if err != nil {
		log.Printf("Failed to save activity to DB: %v", err)
	} else {
		log.Println("Successfully saved activity to DB")
	}
}