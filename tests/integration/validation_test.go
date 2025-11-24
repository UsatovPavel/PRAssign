package integration

import (
	"testing"
)

func TestTeamValidation(t *testing.T) {
	adminToken := genTokenFor("integration-test", true)

	status, _ := postJSONWithToken(t, "/team/add", map[string]interface{}{}, adminToken)
	if status == 201 {
		t.Fatalf("expected validation fail for empty payload")
	}
	status, _ = postJSONWithToken(t, "/team/add", map[string]interface{}{
		"team_name": "x",
		"members":   "not-an-array",
	}, adminToken)
	if status == 201 {
		t.Fatalf("expected validation fail for wrong types")
	}
}

func TestSetIsActiveValidation(t *testing.T) {
	adminToken := genTokenFor("integration-test", true)

	status, _ := postJSONWithToken(t, "/users/setIsActive", map[string]interface{}{}, adminToken)
	if status == 200 {
		t.Fatalf("expected validation error for empty setIsActive")
	}
	status, body := postJSONWithToken(
		t,
		"/users/setIsActive",
		map[string]interface{}{"user_id": "no-such", "is_active": true},
		adminToken,
	)
	if status != 200 && status != 404 {
		t.Fatalf("expected 200 or 404 for setIsActive got %d body=%s", status, string(body))
	}
}

func TestPRCreateValidation(t *testing.T) {
	adminToken := genTokenFor("integration-test", true)

	status, _ := postJSONWithToken(t, "/pullRequest/create", map[string]interface{}{}, adminToken)
	if status == 201 {
		t.Fatalf("expected validation error for empty pr create")
	}
	status, _ = postJSONWithToken(
		t,
		"/pullRequest/create",
		map[string]interface{}{"pull_request_id": "x"},
		adminToken,
	)
	if status == 201 {
		t.Fatalf("expected validation error for missing fields")
	}
}
