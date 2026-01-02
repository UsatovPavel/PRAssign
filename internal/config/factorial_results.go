package config

import (
	"errors"
	"os"
	"strings"
	"time"
)

type FactorialResultsKafkaConfig struct {
	Bootstrap []string
	Topic     string
	Group     string
}

type FactorialRetentionConfig struct {
	TTLSeconds int64
	TimeoutSec int64
}

// Env:
// FACTORIAL_KAFKA_BOOTSTRAP   = "kafka1:9092,kafka2:9092,kafka3:9092"
// FACTORIAL_KAFKA_TOPIC_RESULTS = "factorial.results"
// FACTORIAL_KAFKA_GROUP_RESULTS = "factorial-results-consumer"
func LoadFactorialResultsKafkaConfig() (FactorialResultsKafkaConfig, error) {
	bs := strings.TrimSpace(os.Getenv("FACTORIAL_KAFKA_BOOTSTRAP"))
	topic := strings.TrimSpace(os.Getenv("FACTORIAL_KAFKA_TOPIC_RESULTS"))
	group := strings.TrimSpace(os.Getenv("FACTORIAL_KAFKA_GROUP_RESULTS"))

	if bs == "" {
		return FactorialResultsKafkaConfig{}, errors.New("FACTORIAL_KAFKA_BOOTSTRAP is required")
	}
	if topic == "" {
		return FactorialResultsKafkaConfig{}, errors.New("FACTORIAL_KAFKA_TOPIC_RESULTS is required")
	}
	if group == "" {
		return FactorialResultsKafkaConfig{}, errors.New("FACTORIAL_KAFKA_GROUP_RESULTS is required")
	}

	parts := strings.Split(bs, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return FactorialResultsKafkaConfig{}, errors.New("FACTORIAL_KAFKA_BOOTSTRAP parsed to empty list")
	}

	return FactorialResultsKafkaConfig{
		Bootstrap: out,
		Topic:     topic,
		Group:     group,
	}, nil
}

// FACTORIAL_RESULTS_TTL (duration, e.g. "1h") controls retention delete.
func LoadFactorialRetentionConfig() (FactorialRetentionConfig, error) {
	ttlRaw := strings.TrimSpace(os.Getenv("FACTORIAL_RESULTS_TTL"))
	d, err := time.ParseDuration(ttlRaw)
	if err != nil {
		return FactorialRetentionConfig{}, err
	}
	timeoutRaw := strings.TrimSpace(os.Getenv("FACTORIAL_RESULTS_TIMEOUT"))
	tout, err := time.ParseDuration(timeoutRaw)
	if err != nil {
		return FactorialRetentionConfig{}, err
	}

	return FactorialRetentionConfig{
		TTLSeconds: int64(d.Seconds()),
		TimeoutSec: int64(tout.Seconds()),
	}, nil
}
