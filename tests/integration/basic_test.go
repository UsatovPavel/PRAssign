package integration

import (
	"encoding/json"
	"testing"

	"github.com/UsatovPavel/PRAssign/internal/models"
)

type teamCreateResp struct {
	Team models.Team `json:"team"`
}

type prResp struct {
	Pr models.PullRequest `json:"pr"`
}

func TestTeamCRUD(t *testing.T) {
	adminToken := genTokenFor("integration-test", true)

	tn := unique("team")
	u1 := unique("u1")
	u2 := unique("u2")
	team := models.Team{
		TeamName: tn,
		Members: []models.TeamMember{
			{UserID: u1, Username: "Alice", IsActive: true},
			{UserID: u2, Username: "Bob", IsActive: true},
		},
	}
	status, body := postJSONWithToken(t, "/team/add", team, adminToken)
	if status != 201 {
		t.Fatalf("team add expected 201, got %d body=%s", status, string(body))
	}
	var cre teamCreateResp
	if err := json.Unmarshal(body, &cre); err != nil {
		t.Fatalf("unmarshal create resp: %v", err)
	}
	if cre.Team.TeamName != tn {
		t.Fatalf("team name mismatch")
	}
	status, body = getJSONWithToken(t, "/team/get", map[string]string{"team_name": tn}, adminToken)
	if status != 200 {
		t.Fatalf("team get expected 200 got %d body=%s", status, string(body))
	}
	var got models.Team
	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("unmarshal get: %v", err)
	}
	if got.TeamName != tn {
		t.Fatalf("team get name mismatch")
	}
}

func TestDuplicateTeam(t *testing.T) {
	adminToken := genTokenFor("integration-test", true)

	tn := unique("dupteam")
	u1 := unique("u1")
	team := models.Team{
		TeamName: tn,
		Members:  []models.TeamMember{{UserID: u1, Username: "Alice", IsActive: true}},
	}
	status, _ := postJSONWithToken(t, "/team/add", team, adminToken)
	if status != 201 {
		t.Fatalf("first create expected 201 got %d", status)
	}
	status, body := postJSONWithToken(t, "/team/add", team, adminToken)
	if status == 201 {
		t.Fatalf("duplicate create should not return 201")
	}
	var er struct {
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	_ = json.Unmarshal(body, &er)
	if er.Error.Code == string(models.TeamExists) || status == 400 {
		return
	}
}

func TestPRWorkflowAndMergeIdempotency(t *testing.T) {
	tn := unique("teampr")
	author := unique("author")
	rev1 := unique("rev1")
	rev2 := unique("rev2")
	adminToken := genTokenFor("integration-test", true)
	team := models.Team{
		TeamName: tn,
		Members: []models.TeamMember{
			{UserID: author, Username: "Author", IsActive: true},
			{UserID: rev1, Username: "Rev1", IsActive: true},
			{UserID: rev2, Username: "Rev2", IsActive: true},
		},
	}
	status, body := postJSONWithToken(t, "/team/add", team, adminToken)
	if status != 201 {
		t.Fatalf("team add for pr expected 201 got %d body=%s", status, string(body))
	}
	prID := unique("pr")
	token := genTokenFor(author, false)
	status, body = postJSONWithToken(t, "/pullRequest/create", map[string]string{
		"pull_request_id":   prID,
		"pull_request_name": "Add feature",
		"author_id":         author,
	}, token)
	if status != 201 {
		t.Fatalf("create pr expected 201 got %d body=%s", status, string(body))
	}
	var prr prResp
	if err := json.Unmarshal(body, &prr); err != nil {
		t.Fatalf("unmarshal pr create: %v", err)
	}
	if string(prr.Pr.AuthorID) != author && prr.Pr.AuthorID != author {
		t.Fatalf("pr author mismatch")
	}
	for _, a := range prr.Pr.AssignedReviewers {
		if a == author {
			t.Fatalf("author assigned as reviewer")
		}
	}
	token = genTokenFor(author, false)
	status, body = postJSONWithToken(t, "/pullRequest/merge", map[string]string{"pull_request_id": prID}, token)
	if status != 200 {
		t.Fatalf("merge expected 200 got %d body=%s", status, string(body))
	}
	var m prResp
	if err := json.Unmarshal(body, &m); err != nil {
		t.Fatalf("unmarshal merge: %v", err)
	}
	if m.Pr.Status != models.PRStatusMerged {
		t.Fatalf("merge status expected MERGED, got %s", m.Pr.Status)
	}
	token = genTokenFor(author, false)
	status, body = postJSONWithToken(t, "/pullRequest/merge", map[string]string{"pull_request_id": prID}, token)
	if status != 200 {
		t.Fatalf("merge idempotent expected 200 got %d body=%s", status, string(body))
	}
}

func TestUsersGetReview(t *testing.T) {
	tn := unique("teamrev")
	author := unique("author2")
	rev := unique("rev")
	adminToken := genTokenFor("integration-test", true)
	team := models.Team{
		TeamName: tn,
		Members: []models.TeamMember{
			{UserID: author, Username: "Author", IsActive: true},
			{UserID: rev, Username: "Rev", IsActive: true},
		},
	}
	status, body := postJSONWithToken(t, "/team/add", team, adminToken)
	if status != 201 {
		t.Fatalf("team add expected 201 got %d body=%s", status, string(body))
	}
	prID := unique("pr2")
	token := genTokenFor(author, false)
	status, body = postJSONWithToken(t, "/pullRequest/create", map[string]string{
		"pull_request_id":   prID,
		"pull_request_name": "Do thing",
		"author_id":         author,
	}, token)
	if status != 201 {
		t.Fatalf("create pr expected 201 got %d body=%s", status, string(body))
	}
	token = genTokenFor(rev, false)
	status, body = getJSONWithToken(t, "/users/getReview", map[string]string{"user_id": rev}, token)
	if status != 200 {
		t.Fatalf("getReview expected 200 got %d body=%s", status, string(body))
	}
	var resp struct {
		UserID       string                  `json:"user_id"`
		PullRequests []models.ReviewResponse `json:"pull_requests"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("unmarshal getReview: %v", err)
	}
	if resp.UserID != rev {
		t.Fatalf("getReview user mismatch")
	}
}

func TestReassignEdgeCases(t *testing.T) {
	tn := unique("teamr")
	author := unique("ar")
	a1 := unique("a1")
	adminToken := genTokenFor("integration-test", true)
	team := models.Team{
		TeamName: tn,
		Members: []models.TeamMember{
			{UserID: author, Username: "Author", IsActive: true},
			{UserID: a1, Username: "A1", IsActive: true},
		},
	}
	status, body := postJSONWithToken(t, "/team/add", team, adminToken)
	if status != 201 {
		t.Fatalf("team add expected 201 got %d body=%s", status, string(body))
	}
	prID := unique("prr")
	token := genTokenFor(author, false)
	status, body = postJSONWithToken(t, "/pullRequest/create", map[string]string{
		"pull_request_id":   prID,
		"pull_request_name": "Edge",
		"author_id":         author,
	}, token)
	if status != 201 {
		t.Fatalf("create pr expected 201 got %d body=%s", status, string(body))
	}
	token = genTokenFor(author, false)
	status, _ = postJSONWithToken(t, "/pullRequest/reassign", map[string]string{
		"pull_request_id": prID,
		"old_user_id":     "not-assigned",
	}, token)
	if status == 200 {
		t.Fatalf("reassign should not succeed for not assigned user")
	}
	token = genTokenFor(author, false)
	status, _ = postJSONWithToken(t, "/pullRequest/merge", map[string]string{"pull_request_id": prID}, token)
	if status != 200 {
		t.Fatalf("merge expected 200 got %d", status)
	}
	token = genTokenFor(author, false)
	status, _ = postJSONWithToken(t, "/pullRequest/reassign", map[string]string{
		"pull_request_id": prID,
		"old_user_id":     a1,
	}, token)
	if status == 200 {
		t.Fatalf("reassign on merged PR should not succeed")
	}
}

func BenchmarkCreateTeam(b *testing.B) {
	for i := 0; i < b.N; i++ {
		tn := unique("bench")
		u1 := unique("u1")
		team := models.Team{
			TeamName: tn,
			Members:  []models.TeamMember{{UserID: u1, Username: "Alice", IsActive: true}},
		}
		adminToken := genTokenFor("integration-test", true)
		status, _ := postJSONWithToken(&testing.T{}, "/team/add", team, adminToken)
		if status != 201 {
			b.Fatalf("bench create expected 201 got %d", status)
		}
	}
}
