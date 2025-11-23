package integration

import (
	"time"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"sync"
	"testing"
)

func postJSONNoFatal(url string, payload interface{}) (int, []byte, error) {
	b, err := json.Marshal(payload)
	if err != nil {
		return 0, nil, err
	}
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return resp.StatusCode, body, nil
}

func TestConcurrentTeamCreate(t *testing.T) {
	base := os.Getenv("API_BASE_URL")
	if base == "" {
		base = "http://app_test:8080"
	}
	url := base + "/team/add"
	teamName := unique("concurrent-team")
	payload := map[string]interface{}{
		"team_name": teamName,
		"members": []map[string]interface{}{
			{"user_id": unique("u"), "username": "u", "is_active": true},
		},
	}

	const goroutines = 10
	var wg sync.WaitGroup
	wg.Add(goroutines)

	results := make(chan int, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			status, _, err := postJSONNoFatal(url, payload)
			if err != nil {
				results <- -1
				return
			}
			results <- status
		}()
	}
	wg.Wait()
	close(results)

	successes := 0
	conflicts := 0
	others := 0
	for s := range results {
		if s == 201 || s == 200 {
			successes++
		} else if s == 400 || s == 409 {
			conflicts++
		} else {
			others++
		}
	}
	if successes != 1 {
		t.Fatalf("expected exactly 1 success (201/200), got %d success, %d conflicts, %d others", successes, conflicts, others)
	}
}

func TestConcurrentPRCreateSameID(t *testing.T) {
	base := os.Getenv("API_BASE_URL")
	if base == "" {
		base = "http://app_test:8080"
	}
	url := base + "/pullRequest/create"
	prID := unique("concurrent-pr")
	payload := map[string]interface{}{
		"pull_request_id":   prID,
		"pull_request_name": "concurrent",
		"author_id":         unique("author"),
	}

	const goroutines = 8
	var wg sync.WaitGroup
	wg.Add(goroutines)

	results := make(chan int, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			status, _, err := postJSONNoFatal(url, payload)
			if err != nil {
				results <- -1
				return
			}
			results <- status
		}()
	}
	wg.Wait()
	close(results)

	successes := 0
	conflicts := 0
	others := 0
	for s := range results {
		if s == 201 || s == 200 {
			successes++
		} else if s == 400 || s == 409 {
			conflicts++
		} else {
			others++
		}
	}
	if successes != 1 {
		t.Fatalf("expected exactly 1 success (201/200) for concurrent PR create, got %d success, %d conflicts, %d others", successes, conflicts, others)
	}
}
