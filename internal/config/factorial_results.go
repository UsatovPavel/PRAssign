package config

import (
	"errors"
	"os"
	"strings"
)

type FactorialResultsKafkaConfig struct {
	Bootstrap []string
	Topic     string
	Group     string
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
