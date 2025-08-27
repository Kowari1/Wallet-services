package kafka

import (
	"context"
	"encoding/json"
	"gw-currency-wallet/internal/pkg/logger"
	"time"

	"github.com/segmentio/kafka-go"
)

type ProducerInterface interface {
	Publish(ctx context.Context, key string, value interface{}) error
}

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(broker, topic string, batchSize int, batchTimeout time.Duration) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(broker),
			Topic:                  topic,
			Balancer:               &kafka.LeastBytes{},
			BatchSize:              batchSize,
			BatchTimeout:           batchTimeout,
			RequiredAcks:           kafka.RequireAll,
			Async:                  false,
			AllowAutoTopicCreation: false,
		},
	}
}

func (p *Producer) Publish(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: data,
		Time:  time.Now(),
	}

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		logger.L.Warnw("Kafka publish error", "error", err.Error())
		return err
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
