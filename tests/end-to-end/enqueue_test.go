package endtoend

import (
	"bytes"
	"context"
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

	jobID := enqueueFactorial(t, client, token, 6)
	if jobID == "" {
		t.Fatalf("job_id is empty")
	}
}

func enqueueFactorial(t *testing.T, client http.Client, token string, n int) (jobID string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	jobID = strconv.FormatInt(time.Now().UnixNano(), 10)
	reqBody := []byte(`{"numbers":[` + strconv.Itoa(n) + `]}`)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL()+"/factorial", bytes.NewReader(reqBody))
	if err != nil {
		t.Fatalf("new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Job-Id", jobID)
	if token != "" {
		req.Header.Set("token", token)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("post factorial failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("factorial status = %d, body=%s", resp.StatusCode, string(body))
	}
	return jobID
}
