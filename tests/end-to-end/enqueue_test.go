package endtoend

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"
)

// Simple happy-path for factorial enqueue (no result check).
func TestFactorialEnqueue(t *testing.T) {
	client := http.Client{Timeout: 10 * time.Second}
	token := genToken("e2e-user", false)

	jobID, count := enqueueFactorial(t, client, token, 6)
	if jobID == "" {
		t.Fatalf("job_id is empty")
	}
	if count != 1 {
		t.Fatalf("count expected 1, got %d", count)
	}
}

func enqueueFactorial(t *testing.T, client http.Client, token string, n int) (jobID string, count int) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	reqBody := []byte(`{"numbers":[` + strconv.Itoa(n) + `]}`)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL()+"/factorial", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("token", token)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("post factorial failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		t.Fatalf("factorial status = %d, body=%s", resp.StatusCode, string(body))
	}

	var got struct {
		JobID string `json:"job_id"`
		Count int    `json:"count"`
	}
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("decode resp: %v, body=%s", err, string(body))
	}
	return got.JobID, got.Count
}
