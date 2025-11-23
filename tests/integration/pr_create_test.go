package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

func TestTeamAdd(t *testing.T) {
	time.Sleep(2 * time.Second)

	body := map[string]interface{}{
		"team_name": unique("teamname-testteamadd"),
		"members": []map[string]interface{}{
			{"user_id": unique("u-test-1"), "username": "TestUser", "is_active": true},
		},
	}

	b, _ := json.Marshal(body)

	resp, err := http.Post("http://app_test:8080/team/add", "application/json", bytes.NewBuffer(b))
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 && resp.StatusCode != 200 {
		t.Fatalf("expected 201 or 200, got %d", resp.StatusCode)
	}
}
