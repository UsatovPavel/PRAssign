package integration

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"

	"github.com/UsatovPavel/PRAssign/internal/service"
)

// Smoke: produce factorial tasks to Kafka and ensure they are readable.
// Note: warmUpOffsets is kept to stabilize metadata/leader readiness right after
// topic creation in the test compose stack; plain sleeps were flaky.
func TestFactorialKafkaSmoke(t *testing.T) {
	cfg := loadKafkaTestCfg(t)
	producer := newFactorialProducer(t, cfg.bs, cfg.topic)
	defer producer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err := producer.ProduceTasks(ctx, service.FactorialRequest{
		JobID: cfg.jobID,
		Nums:  cfg.nums,
	})
	if err != nil {
		t.Fatalf("produce: %v", err)
	}

	consumer := newSaramaConsumer(t, cfg.bs)
	defer consumer.Close()

	partitions, err := consumer.Partitions(cfg.topic)
	if err != nil {
		t.Fatalf("partitions: %v", err)
	}
	t.Logf("partitions=%v expected_messages=%d", partitions, len(cfg.nums))
	warmUpOffsets(cfg.bs, cfg.topic, partitions, 5, 500*time.Millisecond)

	deadline := time.Now().Add(30 * time.Second)
	found := consumeJobMessages(t, consumer, cfg.topic, cfg.jobID, partitions, deadline, len(cfg.nums))

	if found < len(cfg.nums) {
		t.Fatalf("expected %d messages for jobId, got %d", len(cfg.nums), found)
	}
}

type kafkaTestCfg struct {
	bs    []string
	topic string
	jobID string
	nums  []int
}

func loadKafkaTestCfg(t *testing.T) kafkaTestCfg {
	t.Helper()

	bootstrap := os.Getenv("FACTORIAL_KAFKA_BOOTSTRAP")
	topic := os.Getenv("FACTORIAL_KAFKA_TOPIC_TASKS")
	if bootstrap == "" || topic == "" {
		t.Skip("FACTORIAL_KAFKA_BOOTSTRAP/FACTORIAL_KAFKA_TOPIC_TASKS not set")
	}

	bs := splitAndTrim(bootstrap)
	if len(bs) == 0 {
		t.Fatal("empty bootstrap after parsing")
	}

	jobID := uuid.NewString()
	nums := []int{5, 6}
	t.Logf("bootstrap=%v topic=%s jobId=%s", bs, topic, jobID)

	return kafkaTestCfg{
		bs:    bs,
		topic: topic,
		jobID: jobID,
		nums:  nums,
	}
}

func newFactorialProducer(t *testing.T, bs []string, topic string) *service.FactorialService {
	t.Helper()

	svc, err := service.NewFactorialService(service.FactorialConfig{
		BootstrapServers: bs,
		TopicTasks:       topic,
	}, nil)
	if err != nil {
		t.Fatalf("producer init: %v", err)
	}
	return svc
}

func newSaramaConsumer(t *testing.T, bs []string) sarama.Consumer {
	t.Helper()

	config := sarama.NewConfig()
	config.Version = sarama.V3_6_0_0
	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer(bs, config)
	if err != nil {
		t.Fatalf("consumer init: %v", err)
	}
	return consumer
}

func splitAndTrim(s string) []string {
	raw := strings.Split(s, ",")
	out := make([]string, 0, len(raw))
	for _, v := range raw {
		v = strings.TrimSpace(v)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

func consumeJobMessages(
	t *testing.T,
	consumer sarama.Consumer,
	topic, jobID string,
	partitions []int32,
	deadline time.Time,
	expected int,
) int {
	found := 0

	for _, p := range partitions {
		pc, err := consumer.ConsumePartition(topic, p, sarama.OffsetOldest)
		if err != nil {
			t.Logf("partition %d: consume err: %v", p, err)
			continue
		}

		n := consumePartition(t, pc, p, jobID, deadline, expected-found)
		pc.Close()

		found += n
		if found >= expected {
			return found
		}
	}

	return found
}

func consumePartition(
	t *testing.T,
	pc sarama.PartitionConsumer,
	partition int32,
	jobID string,
	deadline time.Time,
	expected int,
) int {
	found := 0

	for {
		if time.Now().After(deadline) {
			return found
		}

		select {
		case msg := <-pc.Messages():
			if handleMessage(t, msg, partition, jobID) {
				found++
				if found >= expected {
					return found
				}
			}

		case err := <-pc.Errors():
			if err != nil {
				t.Logf("partition %d: consumer err: %v", partition, err)
			}

		default:
			time.Sleep(50 * time.Millisecond)
		}
	}
}

// warmUpOffsets fetches earliest/latest offsets a few times to ensure metadata/leader ready.
func warmUpOffsets(bs []string, topic string, partitions []int32, retries int, delay time.Duration) {
	clientCfg := sarama.NewConfig()
	clientCfg.Version = sarama.V3_6_0_0

	for i := 0; i < retries; i++ {
		client, err := sarama.NewClient(bs, clientCfg)
		if err == nil {
			for _, p := range partitions {
				_, _ = client.GetOffset(topic, p, sarama.OffsetOldest)
				_, _ = client.GetOffset(topic, p, sarama.OffsetNewest)
			}
			client.Close()
		}
		time.Sleep(delay)
	}
}

func handleMessage(
	t *testing.T,
	msg *sarama.ConsumerMessage,
	partition int32,
	jobID string,
) bool {
	if msg == nil {
		return false
	}

	var item map[string]interface{}
	if err := json.Unmarshal(msg.Value, &item); err != nil {
		t.Logf("partition %d: unmarshal err: %v", partition, err)
		return false
	}

	val, ok := item["jobId"].(string)
	if !ok || val != jobID {
		return false
	}

	t.Logf(
		"got message jobId=%s partition=%d offset=%d",
		jobID,
		partition,
		msg.Offset,
	)

	return true
}
