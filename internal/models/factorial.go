package models

// TaskItem is the payload for factorial requests sent to Kafka.
type TaskItem struct {
	JobID  string `json:"jobId"`
	ItemID int64  `json:"itemId"`
	Input  int    `json:"input"`
}

// ResultItem is produced by Scala into factorial.results.
// output is omitted on failure, error is omitted on success.
type ResultItem struct {
	JobID  string `json:"jobId"`
	ItemID int64  `json:"itemId"`
	Input  int    `json:"input"`
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}
