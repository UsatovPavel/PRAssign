package endtoend

import (
	"context"
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL()+"/health", nil)
	if err != nil {
		t.Fatalf("health request build failed: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("health request failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("health status = %d, body=%s", resp.StatusCode, string(body))
	}
}
