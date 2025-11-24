package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func unique(prefix string) string {
	id := time.Now().UnixNano()
	return prefix + idString(id)
}

func idString(i int64) string {
	return "-" + strconv.FormatInt(i, 10)
}

func genToken() string {
	secret := os.Getenv("AUTH_KEY")
	if secret == "" {
		return ""
	}
	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"user_id":  "integration-test",
		"is_admin": true,
		"iat":      now.Unix(),
		"exp":      now.Add(24 * time.Hour).Unix(),
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := tkn.SignedString([]byte(secret))
	if err != nil {
		return ""
	}
	return signed
}

func TestTeamAdd(t *testing.T) {
	time.Sleep(2 * time.Second)

	body := map[string]interface{}{
		"team_name": unique("teamname-testteam"),
		"members": []map[string]interface{}{
			{"user_id": unique("u-test-1"), "username": "TestUser", "is_active": true},
		},
	}

	b, _ := json.Marshal(body)

	ctx := context.Background()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"http://app_test:8080/team/add",
		bytes.NewBuffer(b),
	)
	if err != nil {
		t.Fatalf("request error: %v", err)
	}
	if err != nil {
		t.Fatalf("request new failed: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token := genToken(); token != "" {
		req.Header.Set("token", token)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 201 or 200, got %d", resp.StatusCode)
	}
}
