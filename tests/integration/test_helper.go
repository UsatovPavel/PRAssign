package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"
)

var baseURL = getBaseURL()

type PRShort struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

func getBaseURL() string {
	v := os.Getenv("API_BASE_URL")
	if v == "" {
		return "http://app_test:8080"
	}
	return v
}

func postJSON(t *testing.T, path string, payload interface{}) (int, []byte) {
	t.Helper()
	b, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	req, err := http.NewRequest("POST", baseURL+path, bytes.NewReader(b))
	if err != nil {
		t.Fatalf("new req: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 10 * time.Second}
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

func getJSON(t *testing.T, path string, params map[string]string) (int, []byte) {
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
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(u.String())
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
	return prefix + "-" + strconv.FormatInt(time.Now().UnixNano(), 10)
}

func nilTesting(_ *testing.B) *testing.T {
	return &testing.T{}
}
