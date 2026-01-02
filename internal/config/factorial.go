package config

import (
	"errors"
	"os"
	"strings"
)

type FactorialKafkaConfig struct {
	Bootstrap  []string
	TopicTasks string
}

// Env:
// FACTORIAL_KAFKA_BOOTSTRAP = "kafka1:9092,kafka2:9092,kafka3:9092"
// FACTORIAL_KAFKA_TOPIC_TASKS = "factorial.tasks"
func LoadFactorialKafkaConfig() (FactorialKafkaConfig, error) {
	bs := strings.TrimSpace(os.Getenv("FACTORIAL_KAFKA_BOOTSTRAP"))
	topic := strings.TrimSpace(os.Getenv("FACTORIAL_KAFKA_TOPIC_TASKS"))

	if bs == "" {
		return FactorialKafkaConfig{}, errors.New("FACTORIAL_KAFKA_BOOTSTRAP is required")
	}
	if topic == "" {
		return FactorialKafkaConfig{}, errors.New("FACTORIAL_KAFKA_TOPIC_TASKS is required")
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
		return FactorialKafkaConfig{}, errors.New("FACTORIAL_KAFKA_BOOTSTRAP parsed to empty list")
	}

	return FactorialKafkaConfig{
		Bootstrap:  out,
		TopicTasks: topic,
	}, nil
}
