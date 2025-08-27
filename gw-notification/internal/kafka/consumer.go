package kafka

import (
	"context"
	"encoding/json"
	"gw-notification/internal/models"
	"gw-notification/internal/pkg/logger"
	"gw-notification/internal/repository"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader     *kafka.Reader
	repository *repository.MongoRepository
}

func NewConsumer(broker, topic, groupID string, storage *repository.MongoRepository) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  []string{broker},
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3,
			MaxBytes: 10e6,
		}),
		repository: storage,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	logger.L.Info("Kafka consumer started")

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}

			logger.L.Errorw("Kafka read error", "error", err.Error())
			continue
		}

		var evt models.EventMessage
		if err := json.Unmarshal(m.Value, &evt); err != nil {
			logger.L.Warnw("Invalid event JSON", "value", string(m.Value), "error", err.Error())
			continue
		}

		if err := c.repository.SaveEvents(ctx, evt); err != nil {
			logger.L.Errorw("Mongo save error", "event_id", evt.EventID, "error", err.Error())
			continue
		}

		logger.L.Infow("Event saved", "event_id", evt.EventID, "user_id", evt.UserID)
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
