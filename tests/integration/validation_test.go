package integration

import (
	"testing"
)

func TestTeamValidation(t *testing.T) {
	status, _ := postJSON(t, "/team/add", map[string]interface{}{})
	if status == 201 {
		t.Fatalf("expected validation fail for empty payload")
	}
	status, _ = postJSON(t, "/team/add", map[string]interface{}{
		"team_name": "x",
		"members":   "not-an-array",
	})
	if status == 201 {
		t.Fatalf("expected validation fail for wrong types")
	}
}

func TestSetIsActiveValidation(t *testing.T) {
	status, _ := postJSON(t, "/users/setIsActive", map[string]interface{}{})
	if status == 200 {
		t.Fatalf("expected validation error for empty setIsActive")
	}
	status, body := postJSON(t, "/users/setIsActive", map[string]interface{}{"user_id": "no-such", "is_active": true})
	if status != 200 && status != 404 {
		t.Fatalf("expected 200 or 404 for setIsActive got %d body=%s", status, string(body))
	}
}

func TestPRCreateValidation(t *testing.T) {
	status, _ := postJSON(t, "/pullRequest/create", map[string]interface{}{})
	if status == 201 {
		t.Fatalf("expected validation error for empty pr create")
	}
	status, _ = postJSON(t, "/pullRequest/create", map[string]interface{}{"pull_request_id": "x"})
	if status == 201 {
		t.Fatalf("expected validation error for missing fields")
	}
}
