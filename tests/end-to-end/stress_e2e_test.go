package endtoend

import (
	"math/rand"
	"net/http"
	"testing"
	"time"
)

// Stress: 100 items of ~200 (+rand10) to check throughput and aggregation.
func TestFactorialBulkLargeNumbers(t *testing.T) {
	client := http.Client{Timeout: 30 * time.Second}
	token := genToken("e2e-user", false)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	nums := make([]int, 100)
	for i := range nums {
		nums[i] = 200 + r.Intn(10) // ~200 with small spread
	}

	jobID := enqueueFactorialList(t, client, token, nums)

	result := waitForFactorialResultWithTimeout(t, client, token, jobID, 180*time.Second)

	if result.TotalItems != 100 {
		t.Fatalf("total_items expected 100 got %d", result.TotalItems)
	}
	if result.FailedItems > 0 {
		t.Fatalf("factorial failed: failed_items=%d", result.FailedItems)
	}
	if len(result.Items) != 100 {
		t.Fatalf("expected 100 result items, got %d", len(result.Items))
	}
	missing := 0
	for _, it := range result.Items {
		if it.Output == nil {
			missing++
		}
	}
	if missing > 0 {
		t.Fatalf("missing outputs for %d items", missing)
	}
}

// Stress: 100 jobs, каждый по ~100 чисел (~200 + rand10).
func TestFactorialHundredJobsHundredItems(t *testing.T) {
	client := http.Client{Timeout: 30 * time.Second}
	token := genToken("e2e-user", false)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	const jobs = 100
	const itemsPerJob = 100

	jobIDs := make([]string, 0, jobs)
	for j := 0; j < jobs; j++ {
		nums := make([]int, itemsPerJob)
		for i := range nums {
			nums[i] = 200 + r.Intn(10)
		}
		jobID := enqueueFactorialList(t, client, token, nums)
		jobIDs = append(jobIDs, jobID)
	}

	deadline := time.Now().Add(300 * time.Second)
	results := make(map[string]FactorialResult, jobs)

	for len(results) < jobs {
		if time.Now().After(deadline) {
			t.Fatalf("timeout waiting for factorial results: got %d/%d", len(results), jobs)
		}
		for _, jobID := range jobIDs {
			if _, ok := results[jobID]; ok {
				continue
			}
			res := getFactorialResult(t, client, token, jobID)
			if res.TotalItems > 0 && res.DoneItems+res.FailedItems >= res.TotalItems {
				results[jobID] = res
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	for _, res := range results {
		if res.TotalItems != itemsPerJob {
			t.Fatalf("job %s total_items expected %d got %d", res.Status, itemsPerJob, res.TotalItems)
		}
		if res.FailedItems > 0 {
			t.Fatalf("job %s failed_items=%d", res.Status, res.FailedItems)
		}
		if len(res.Items) != itemsPerJob {
			t.Fatalf("job %s expected %d result items, got %d", res.Status, itemsPerJob, len(res.Items))
		}
		missing := 0
		for _, it := range res.Items {
			if it.Output == nil {
				missing++
			}
		}
		if missing > 0 {
			t.Fatalf("job %s missing outputs for %d items", res.Status, missing)
		}
	}
}
