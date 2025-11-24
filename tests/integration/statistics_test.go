package integration

import (
	"encoding/json"
	"testing"
)

type prCreateResp struct {
	Pr struct {
		PullRequestID     string   `json:"pull_request_id"`
		AssignedReviewers []string `json:"assigned_reviewers"`
	} `json:"pr"`
}

func fetchPRsStats(t *testing.T, adminToken string) map[string]int {
	t.Helper()
	status, body := getJSONWithToken(t, "/statistics/assignments/pullrequests", nil, adminToken)
	if status != 200 {
		t.Fatalf("fetch prs stats expected 200 got %d body=%s", status, string(body))
	}
	var resp struct {
		ByPR map[string]int `json:"assignments_by_pr"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("unmarshal prs stats: %v", err)
	}
	return resp.ByPR
}

func fetchUsersStats(t *testing.T, adminToken string) map[string]int {
	t.Helper()
	status, body := getJSONWithToken(t, "/statistics/assignments/users", nil, adminToken)
	if status != 200 {
		t.Fatalf("fetch users stats expected 200 got %d body=%s", status, string(body))
	}
	var resp struct {
		ByUser map[string]int `json:"assignments_by_user"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("unmarshal users stats: %v", err)
	}
	return resp.ByUser
}

func createPRAndGetAssigned(t *testing.T, author string) (string, int) {
	t.Helper()
	prID := unique("pr")
	status, body := postJSONWithToken(t,
		"/pullRequest/create",
		map[string]string{
			"pull_request_id":   prID,
			"pull_request_name": "Feature",
			"author_id":         author,
		},
		genTokenFor(author, false),
	)
	if status != 201 {
		t.Fatalf("create pr expected 201 got %d body=%s", status, string(body))
	}
	var r prCreateResp
	if err := json.Unmarshal(body, &r); err != nil {
		t.Fatalf("unmarshal pr: %v", err)
	}
	return r.Pr.PullRequestID, len(r.Pr.AssignedReviewers)
}

func TestStatsAssignmentsByPRsAndUsers(t *testing.T) {
	_, author, rev1, rev2 := setupPRTeam(t)

	adminToken := genTokenFor("integration-test", true)

	beforePRs := fetchPRsStats(t, adminToken)
	beforeUsers := fetchUsersStats(t, adminToken)

	pr1, cnt1 := createPRAndGetAssigned(t, author)
	pr2, cnt2 := createPRAndGetAssigned(t, author)

	afterPRs := fetchPRsStats(t, adminToken)
	afterUsers := fetchUsersStats(t, adminToken)

	if afterPRs[pr1]-beforePRs[pr1] != cnt1 {
		t.Fatalf("pr1 delta mismatch got %d want %d", afterPRs[pr1]-beforePRs[pr1], cnt1)
	}
	if afterPRs[pr2]-beforePRs[pr2] != cnt2 {
		t.Fatalf("pr2 delta mismatch got %d want %d", afterPRs[pr2]-beforePRs[pr2], cnt2)
	}

	sumBefore := 0
	for _, v := range beforeUsers {
		sumBefore += v
	}
	sumAfter := 0
	for _, v := range afterUsers {
		sumAfter += v
	}
	if sumAfter-sumBefore != cnt1+cnt2 {
		t.Fatalf("total assignments delta mismatch got %d want %d", sumAfter-sumBefore, cnt1+cnt2)
	}
	if afterUsers[rev1]-beforeUsers[rev1] == 0 && afterUsers[rev2]-beforeUsers[rev2] == 0 {
		t.Fatalf("expected at least one assignment delta for reviewers; before=%v after=%v", beforeUsers, afterUsers)
	}
}

func fetchUserCount(t *testing.T, adminToken, userID string) int {
	t.Helper()
	status, body := getJSONWithToken(t, "/statistics/assignments/user/"+userID, nil, adminToken)
	if status != 200 {
		t.Fatalf("fetch user stats expected 200 got %d body=%s", status, string(body))
	}
	var resp struct {
		UserID           string `json:"user_id"`
		AssignmentsCount int    `json:"assignments_count"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("unmarshal user stats: %v", err)
	}
	return resp.AssignmentsCount
}

func TestStatsAssignmentsForUser(t *testing.T) {
	_, author, rev1, rev2 := setupPRTeam(t)

	adminToken := genTokenFor("integration-test", true)

	before1 := fetchUserCount(t, adminToken, rev1)
	before2 := fetchUserCount(t, adminToken, rev2)

	_, _ = createPRAndGetAssigned(t, author)
	_, _ = createPRAndGetAssigned(t, author)

	after1 := fetchUserCount(t, adminToken, rev1)
	after2 := fetchUserCount(t, adminToken, rev2)

	diff1 := after1 - before1
	diff2 := after2 - before2
	if diff1 < 0 || diff1 > 2 {
		t.Fatalf("assignments count delta out of expected range for %s got %d", rev1, diff1)
	}
	if diff2 < 0 || diff2 > 2 {
		t.Fatalf("assignments count delta out of expected range for %s got %d", rev2, diff2)
	}
}
