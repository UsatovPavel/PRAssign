package models

// TaskItem is the payload for factorial requests sent to Kafka.
type TaskItem struct {
	JobID  string `json:"jobId"`
	ItemID int64  `json:"itemId"`
	Input  int    `json:"input"`
}



