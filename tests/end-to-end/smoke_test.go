package endtoend

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/UsatovPavel/PRAssign/internal/config"
	"github.com/golang-jwt/jwt/v4"
)

func baseURL() string {
	if v := os.Getenv("API_BASE_URL"); v != "" {
		return v
	}
	return "http://localhost:8080"
}

func genToken(user string, isAdmin bool) string {
	secret := os.Getenv("AUTH_KEY")
	if secret == "" {
		return ""
	}
	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"user_id":  user,
		"is_admin": isAdmin,
		"iat":      now.Unix(),
		"exp":      now.Add(config.TokenExpiration).Unix(),
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tkn.SignedString([]byte(secret))
	if err != nil {
		return ""
	}
	return signed
}

// Smoke test: service is reachable and health endpoint responds 200.
func TestHealth(t *testing.T) {
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(baseURL() + "/health")
	if err != nil {
		t.Fatalf("health request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("health status = %d, body=%s", resp.StatusCode, string(body))
	}
}

// Simple happy-path for factorial enqueue.
func TestFactorialEnqueue(t *testing.T) {
	client := http.Client{Timeout: 10 * time.Second}
	token := genToken("e2e-user", false)

	reqBody := []byte(`{"numbers":[6]}`)
	req, err := http.NewRequest(http.MethodPost, baseURL()+"/factorial", bytes.NewReader(reqBody))
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
	if got.JobID == "" {
		t.Fatalf("job_id is empty")
	}
	if got.Count != 1 {
		t.Fatalf("count expected 1, got %d", got.Count)
	}
}

