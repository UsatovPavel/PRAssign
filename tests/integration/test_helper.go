package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/UsatovPavel/PRAssign/internal/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var baseURL = getBaseURL()

func getBaseURL() string {
	v := os.Getenv("API_BASE_URL")
	if v == "" {
		return "http://app_test:8080"
	}
	return v
}

func genTokenFor(user string, isAdmin bool) string {
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

func postJSONWithToken(t *testing.T, path string, payload interface{}, token string) (int, []byte) {
	t.Helper()
	b, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.HTTPClientTimeoutLong)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+path, bytes.NewReader(b))
	if err != nil {
		t.Fatalf("new req: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("token", token)
	}

	client := &http.Client{Timeout: config.HTTPClientTimeoutLong}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("do post: %v", err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return resp.StatusCode, buf.Bytes()
}

func getJSONWithToken(t *testing.T, path string, params map[string]string, token string) (int, []byte) {
	t.Helper()

	u, err := url.Parse(baseURL + path)
	if err != nil {
		t.Fatalf("parse url: %v", err)
	}
	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	ctx, cancel := context.WithTimeout(context.Background(), config.HTTPClientTimeoutLong)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		t.Fatalf("new req: %v", err)
	}
	if token != "" {
		req.Header.Set("token", token)
	}

	client := &http.Client{Timeout: config.HTTPClientTimeoutLong}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return resp.StatusCode, buf.Bytes()
}

func unique(prefix string) string {
	return prefix + "-" + uuid.NewString()[:8]
}
