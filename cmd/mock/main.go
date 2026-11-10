package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Result represents a mocked exam result.
type Result struct {
	Subject string `json:"subject"`
	Traz    int    `json:"traz"`
	Status  string `json:"status"`
}

// MazeResponse represents the mocked upstream payload.
type MazeResponse struct {
	Data       []Result `json:"data"`
	TotalPages int      `json:"total_pages"`
	Page       int      `json:"page"`
}

func main() {
	rand.Seed(time.Now().UnixNano())
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/students/{id}/results", func(w http.ResponseWriter, r *http.Request) {
		// Simulate upstream instability: 20% chance of 429 Too Many Requests
		if rand.Float32() < 0.2 {
			log.Printf("Simulating 429 Too Many Requests for student %s", r.PathValue("id"))
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		// Simulate network latency (0 to 1000 milliseconds)
		delay := time.Duration(rand.Intn(1000)) * time.Millisecond
		time.Sleep(delay)

		mockData := MazeResponse{
			Data: []Result{
				{Subject: "ریاضی", Traz: 4500, Status: "scored"},
				{Subject: "فیزیک", Traz: 5800, Status: "scored"},
				{Subject: "شیمی", Traz: 7200, Status: "scored"},
				{Subject: "ادبیات", Traz: 0, Status: "pending"}, // Should be ignored by the client
			},
			TotalPages: 1,
			Page:       1,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockData)
		log.Printf("Successfully returned data with %v delay", delay)
	})

	fmt.Println("🚀 Mock Maze Simulator is running on port :8081...")
	log.Fatal(http.ListenAndServe(":8081", mux))
}
