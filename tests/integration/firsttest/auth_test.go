package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestAuthTokenGeneration(t *testing.T) {
	time.Sleep(1 * time.Second)

	body := map[string]interface{}{
		"username": unique("authuser"),
	}

	b, _ := json.Marshal(body)
	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://app_test:8080/auth/token", bytes.NewBuffer(b))
	if err != nil {
		t.Fatalf("request error: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 200/201, got %d", resp.StatusCode)
	}

	var out struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if out.Token == "" {
		t.Fatalf("expected non-empty token")
	}
}
