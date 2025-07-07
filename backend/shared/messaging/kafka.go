package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/IBM/sarama"
	"korun.io/shared/config"
	"korun.io/shared/events"
)

type Producer interface {
	PublishEvent(ctx context.Context, topic string, event *events.Event) error
	Close() error
}

type KafkaProducer struct {
	producer sarama.SyncProducer
}

func NewKafkaProducer(cfg *config.KafkaConfig) (*KafkaProducer, error) {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Return.Errors = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = cfg.RetryAttempts
	saramaConfig.Producer.Retry.Backoff = cfg.RetryBackoff

	producer, err := sarama.NewSyncProducer(cfg.Brokers, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &KafkaProducer{producer: producer}, nil
}

func (p *KafkaProducer) PublishEvent(ctx context.Context, topic string, event *events.Event) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(event.ID),
		Value: sarama.ByteEncoder(eventBytes),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("event_type"),
				Value: []byte(event.Type),
			},
			{
				Key:   []byte("source"),
				Value: []byte(event.Source),
			},
		},
	}

	partition, offset, err := p.producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %w", err)
	}

	slog.Info("Event published to Kafka", "topic", topic, "partition", partition, "offset", offset, "event_id", event.ID)

	return nil
}

func (p *KafkaProducer) Close() error {
	return p.producer.Close()
}
