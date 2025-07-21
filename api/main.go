package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/segmentio/kafka-go"
)

type Activity struct {
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

var (
	kafkaBroker   = os.Getenv("KAFKA_BROKER")
	kafkaTopic    = os.Getenv("KAFKA_TOPIC")
	clickhouseAddr = os.Getenv("CLICKHOUSE_ADDR")
	kafkaWriter   *kafka.Writer
	db            clickhouse.Conn
)

func main() {
	var err error

	// Initialize Kafka Writer
	kafkaWriter = &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    kafkaTopic,
		Balancer: &kafka.LeastBytes{},
	}
	defer kafkaWriter.Close()

	// Initialize ClickHouse connection
	db, err = connectToClickHouse()
	if err != nil {
		log.Fatalf("Failed to connect to ClickHouse: %v", err)
	}
	log.Println("Successfully connected to ClickHouse")
	defer db.Close()

	// Setup HTTP server endpoints
	http.HandleFunc("/track", trackHandler)
	http.HandleFunc("/activities", getActivitiesHandler)

	log.Println("API server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func connectToClickHouse() (clickhouse.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{clickhouseAddr},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "streamtrack-api", Version: "0.1"},
			},
		},
	})

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}
	return conn, nil
}

func trackHandler(w http.ResponseWriter, r *http.Request) {
	// ... (This handler remains exactly the same as before)
}

func getActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query(context.Background(), "SELECT user_id, action, created_at FROM activities ORDER BY created_at DESC")
	if err != nil {
		http.Error(w, "Failed to query database", http.StatusInternalServerError)
		log.Printf("db.Query error: %v", err)
		return
	}
	defer rows.Close()

	var activities []Activity
	for rows.Next() {
		var activity Activity
		if err := rows.Scan(&activity.UserID, &activity.Action, &activity.CreatedAt); err != nil {
			http.Error(w, "Failed to scan database row", http.StatusInternalServerError)
			log.Printf("rows.Scan error: %v", err)
			return
		}
		activities = append(activities, activity)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(activities); err != nil {
		http.Error(w, "Failed to encode response to JSON", http.StatusInternalServerError)
	}
}