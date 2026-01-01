package endtoend

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestFactorialEndToEnd(t *testing.T) {
	client := http.Client{Timeout: 10 * time.Second}
	token := genToken("e2e-user", false)

	jobID, _ := enqueueFactorial(t, client, token, 6)

	result := waitForFactorialResult(t, client, token, jobID)

	if result.TotalItems != 1 {
		t.Fatalf("total_items expected 1 got %d", result.TotalItems)
	}
	if result.FailedItems > 0 {
		t.Fatalf("factorial failed")
	}
	if len(result.Items) == 0 || result.Items[0].Output == nil || *result.Items[0].Output != "720" {
		t.Fatalf("unexpected result")
	}
}

func waitForFactorialResult(t *testing.T, client http.Client, token, jobID string) FactorialResult {
	t.Helper()
	deadline := time.Now().Add(20 * time.Second)

	for {
		res := getFactorialResult(t, client, token, jobID)
		if res.TotalItems > 0 && res.DoneItems+res.FailedItems >= res.TotalItems {
			return res
		}
		if time.Now().After(deadline) {
			t.Fatalf("timeout waiting for factorial result")
		}
		time.Sleep(300 * time.Millisecond)
	}
}

type FactorialResult struct {
	Status      string `json:"status"`
	TotalItems  int    `json:"total_items"`
	DoneItems   int    `json:"done_items"`
	FailedItems int    `json:"failed_items"`
	Items       []struct {
		Input  int     `json:"input"`
		Status string  `json:"status"`
		Output *string `json:"output"`
		Error  *string `json:"error"`
	} `json:"items"`
}

func getFactorialResult(t *testing.T, client http.Client, token, jobID string) FactorialResult {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL()+"/factorial/"+jobID+"/result", nil)
	if err != nil {
		t.Fatalf("build get result req: %v", err)
	}
	if token != "" {
		req.Header.Set("token", token)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("get result failed: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == http.StatusNotFound {
		return FactorialResult{} // ещё не готово
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("result status = %d body=%s", resp.StatusCode, string(body))
	}

	var res FactorialResult
	if err := json.Unmarshal(body, &res); err != nil {
		t.Fatalf("decode result: %v", err)
	}
	return res
}
