package domain

// EventPublisher defines the interface for publishing account events
type EventPublisher interface {
	PublishAccountCreated(accountID string, status string) error
	PublishAccountStatusChanged(accountID string, status string) error
}
