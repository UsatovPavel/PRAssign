package service

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"

	"github.com/IBM/sarama"

	"github.com/UsatovPavel/PRAssign/internal/models"
)

type FactorialService struct {
	producer   sarama.SyncProducer
	topicTasks string
	l          *slog.Logger
}

type FactorialConfig struct {
	BootstrapServers []string
	TopicTasks       string
}

func NewFactorialService(cfg FactorialConfig, l *slog.Logger) (*FactorialService, error) {
	prodCfg := sarama.NewConfig()
	prodCfg.Producer.Return.Successes = true
	prodCfg.Producer.RequiredAcks = sarama.WaitForAll
	prodCfg.Producer.Idempotent = true
	prodCfg.Net.MaxOpenRequests = 1

	prod, err := sarama.NewSyncProducer(cfg.BootstrapServers, prodCfg)
	if err != nil {
		return nil, err
	}

	return &FactorialService{
		producer:   prod,
		topicTasks: cfg.TopicTasks,
		l:          l,
	}, nil
}

func (s *FactorialService) Close() error {
	return s.producer.Close()
}

type FactorialRequest struct {
	JobID string
	Nums  []int
}

type FactorialResponse struct {
	JobID string `json:"job_id"`
	Count int    `json:"count"`
}

// ProduceTasks publishes one message per input number to Kafka.
func (s *FactorialService) ProduceTasks(ctx context.Context, req FactorialRequest) (FactorialResponse, error) {
	jobID := req.JobID
	if strings.TrimSpace(jobID) == "" {
		return FactorialResponse{}, errors.New("job_id is required")
	}

	for idx, n := range req.Nums {
		select {
		case <-ctx.Done():
			return FactorialResponse{}, ctx.Err()
		default:
		}

		msgBody, err := json.Marshal(models.TaskItem{
			JobID:  jobID,
			ItemID: int64(idx),
			Input:  n,
		})
		if err != nil {
			return FactorialResponse{}, err
		}

		msg := &sarama.ProducerMessage{
			Topic: s.topicTasks,
			Key:   sarama.StringEncoder(jobID),
			Value: sarama.ByteEncoder(msgBody),
		}

		partition, offset, err := s.producer.SendMessage(msg)
		if err != nil {
			if s.l != nil {
				// Sarama may return a ProducerErrors; unwrap to log.
				s.l.Error("factorial produce failed", "err", err, "job_id", jobID, "item_id", idx)
			}
			return FactorialResponse{}, err
		}

		if s.l != nil {
			s.l.Info(
				"factorial produce ok",
				"topic", s.topicTasks,
				"job_id", jobID,
				"item_id", idx,
				"partition", partition,
				"offset", offset,
			)
		}
	}

	return FactorialResponse{JobID: jobID, Count: len(req.Nums)}, nil
}
