package infrastructure

import (
	"context"
	"encoding/json"
	"log"

	"github.com/DavidRodriguez-create/pay-and-go/services/card/domain"
	"github.com/segmentio/kafka-go"
)

// AccountEvent represents an event from the account service
type AccountEvent struct {
	Type      string `json:"type"` // "account.created" or "account.status_changed"
	AccountID string `json:"account_id"`
	Status    string `json:"status"` // "ACTIVE", "BLOCKED", "DELETED"
}

// KafkaAccountConsumer consumes account events from Kafka
type KafkaAccountConsumer struct {
	reader      *kafka.Reader
	accountRepo domain.AccountCacheRepository
	stopChan    chan struct{}
}

// NewKafkaAccountConsumer creates a new Kafka consumer for account events
func NewKafkaAccountConsumer(
	brokers []string,
	topic string,
	groupID string,
	accountRepo domain.AccountCacheRepository,
) *KafkaAccountConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	return &KafkaAccountConsumer{
		reader:      reader,
		accountRepo: accountRepo,
		stopChan:    make(chan struct{}),
	}
}

// Start begins consuming messages from Kafka
func (c *KafkaAccountConsumer) Start(ctx context.Context) error {
	log.Println("Starting Kafka account event consumer...")

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("Context cancelled, stopping consumer...")
				return
			case <-c.stopChan:
				log.Println("Stop signal received, stopping consumer...")
				return
			default:
				msg, err := c.reader.ReadMessage(ctx)
				if err != nil {
					if err == context.Canceled {
						return
					}
					log.Printf("Error reading message: %v\n", err)
					continue
				}

				if err := c.handleMessage(msg); err != nil {
					log.Printf("Error handling message: %v\n", err)
				}
			}
		}
	}()

	return nil
}

// Stop stops the Kafka consumer
func (c *KafkaAccountConsumer) Stop() error {
	close(c.stopChan)
	return c.reader.Close()
}

// handleMessage processes a single Kafka message
func (c *KafkaAccountConsumer) handleMessage(msg kafka.Message) error {
	var event AccountEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return err
	}

	log.Printf("Received account event: type=%s, account_id=%s, status=%s\n",
		event.Type, event.AccountID, event.Status)

	// Convert status string to AccountStatus
	status := domain.AccountStatus(event.Status)

	// Upsert account cache
	accountCache := domain.NewAccountCache(event.AccountID, status)
	if err := c.accountRepo.Upsert(accountCache); err != nil {
		return err
	}

	log.Printf("Updated account cache: account_id=%s, status=%s\n",
		event.AccountID, event.Status)

	return nil
}
