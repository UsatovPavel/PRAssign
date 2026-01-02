package service

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/IBM/sarama"

	"github.com/UsatovPavel/PRAssign/internal/models"
	"github.com/UsatovPavel/PRAssign/internal/repository"
)

type FactorialResultConsumerConfig struct {
	Bootstrap []string
	Group     string
	Topic     string
}

// FactorialResultConsumer consumes factorial.results and upserts rows into Postgres.
type FactorialResultConsumer struct {
	repo   repository.FactorialRepository
	logger *slog.Logger
	cfg    FactorialResultConsumerConfig
}

func NewFactorialResultConsumer(repo repository.FactorialRepository, cfg FactorialResultConsumerConfig, l *slog.Logger) *FactorialResultConsumer {
	return &FactorialResultConsumer{repo: repo, logger: l, cfg: cfg}
}

func (c *FactorialResultConsumer) Run(ctx context.Context) error {
	config := sarama.NewConfig()
	config.Version = sarama.V3_6_0_0
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewConsumerGroup(c.cfg.Bootstrap, c.cfg.Group, config)
	if err != nil {
		return err
	}
	defer client.Close()

	handler := &factorialResultHandler{repo: c.repo, logger: c.logger}
	for {
		if err := client.Consume(ctx, []string{c.cfg.Topic}, handler); err != nil {
			return err
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

type factorialResultHandler struct {
	repo   repository.FactorialRepository
	logger *slog.Logger
}

func (h *factorialResultHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *factorialResultHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *factorialResultHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var res models.ResultItem
		if err := json.Unmarshal(msg.Value, &res); err != nil {
			if h.logger != nil {
				h.logger.Error("factorial result: unmarshal failed", "err", err)
			}
			session.MarkMessage(msg, "")
			continue
		}

		out := res.Result

		status := "done"
		if res.Error != "" {
			status = "failed"
		}

		row := repository.FactorialResultRow{
			JobID:  res.JobID,
			ItemID: res.ItemID,
			Input:  res.Input,
			Status: status,
		}
		if out != "" {
			row.Output = &out
		}
		if res.Error != "" {
			row.Error = &res.Error
		}

		if err := h.repo.UpsertResult(session.Context(), row); err != nil {
			if h.logger != nil {
				h.logger.Error("factorial result: upsert failed", "err", err, "job_id", res.JobID, "item_id", res.ItemID)
			}
			// do not mark offset to retry
			continue
		}

		session.MarkMessage(msg, "")
	}
	return nil
}
