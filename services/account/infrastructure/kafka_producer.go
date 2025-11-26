package infrastructure

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

// AccountEvent represents an event from the account service
type AccountEvent struct {
	Type      string `json:"type"` // "account.created" or "account.status_changed"
	AccountID string `json:"account_id"`
	Status    string `json:"status"` // "ACTIVE", "BLOCKED", "DELETED"
}

// KafkaProducer handles publishing events to Kafka
type KafkaProducer struct {
	writer *kafka.Writer
}

// NewKafkaProducer creates a new Kafka producer
func NewKafkaProducer(brokers []string, topic string) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaProducer{
		writer: writer,
	}
}

// PublishAccountCreated publishes an account.created event
func (p *KafkaProducer) PublishAccountCreated(accountID string, status string) error {
	event := AccountEvent{
		Type:      "account.created",
		AccountID: accountID,
		Status:    status,
	}

	return p.publish(event)
}

// PublishAccountStatusChanged publishes an account.status_changed event
func (p *KafkaProducer) PublishAccountStatusChanged(accountID string, status string) error {
	event := AccountEvent{
		Type:      "account.status_changed",
		AccountID: accountID,
		Status:    status,
	}

	return p.publish(event)
}

// publish sends an event to Kafka
func (p *KafkaProducer) publish(event AccountEvent) error {
	value, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Value: value,
	}

	err = p.writer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Printf("Failed to publish event: %v\n", err)
		return err
	}

	log.Printf("Published event: type=%s, account_id=%s, status=%s\n",
		event.Type, event.AccountID, event.Status)
	return nil
}

// Close closes the Kafka writer
func (p *KafkaProducer) Close() error {
	return p.writer.Close()
}
