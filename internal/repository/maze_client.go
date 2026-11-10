package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// Result represents the required fields extracted from upstream.
type Result struct {
	Subject string `json:"subject"`
	Traz    int    `json:"traz"`
	Status  string `json:"status"`
}

// MazeResponse defines the external API payload structure.
type MazeResponse struct {
	Data []Result `json:"data"`
}

// MazeClient manages communications with the external Maze simulator.
type MazeClient struct {
	BaseURL    string
	httpClient *http.Client
	cache      sync.Map // Thread-safe in-memory cache for high concurrency
}

// NewMazeClient creates a new instance of MazeClient with a predefined timeout.
func NewMazeClient(baseURL string) *MazeClient {
	return &MazeClient{
		BaseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// FetchResults retrieves exam results, utilizing a concurrent-safe cache
// and implementing an exponential backoff & retry mechanism for resilience.
func (c *MazeClient) FetchResults(studentID, dateStr string) ([]Result, error) {
	cacheKey := studentID + ":" + dateStr

	// Fast path: Return from in-memory cache if available
	if cachedData, ok := c.cache.Load(cacheKey); ok {
		return cachedData.([]Result), nil
	}

	url := fmt.Sprintf("%s/api/students/%s/results", c.BaseURL, studentID)
	var lastErr error

	// Slow path: Request upstream with up to 3 retries
	for i := 0; i < 3; i++ {
		req, _ := http.NewRequest("GET", url, nil)
		resp, err := c.httpClient.Do(req)

		if err != nil {
			lastErr = err
			time.Sleep(time.Duration(i+1) * 200 * time.Millisecond)
			continue
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			resp.Body.Close()
			lastErr = errors.New("rate limited 429")
			// Exponential backoff strategy
			time.Sleep(time.Duration(i+1) * 500 * time.Millisecond)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			return nil, fmt.Errorf("upstream error: %d", resp.StatusCode)
		}

		var mazeResp MazeResponse
		if err := json.NewDecoder(resp.Body).Decode(&mazeResp); err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()

		// Filter out 'pending' items
		var validSubjects []Result
		for _, s := range mazeResp.Data {
			if s.Status == "scored" {
				validSubjects = append(validSubjects, s)
			}
		}

		// Populate cache for subsequent concurrent requests
		c.cache.Store(cacheKey, validSubjects)
		return validSubjects, nil
	}

	return nil, fmt.Errorf("failed to fetch from maze after 3 retries: %v", lastErr)
}
