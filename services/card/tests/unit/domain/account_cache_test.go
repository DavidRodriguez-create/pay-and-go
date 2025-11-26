package domain_test

import (
	"testing"

	"github.com/DavidRodriguez-create/pay-and-go/services/card/domain"
)

func TestNewAccountCache(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		status   domain.AccountStatus
		expected *domain.AccountCache
	}{
		{
			name:   "Create active account cache",
			id:     "acc-123",
			status: domain.AccountStatusActive,
			expected: &domain.AccountCache{
				ID:     "acc-123",
				Status: domain.AccountStatusActive,
			},
		},
		{
			name:   "Create blocked account cache",
			id:     "acc-456",
			status: domain.AccountStatusBlocked,
			expected: &domain.AccountCache{
				ID:     "acc-456",
				Status: domain.AccountStatusBlocked,
			},
		},
		{
			name:   "Create deleted account cache",
			id:     "acc-789",
			status: domain.AccountStatusDeleted,
			expected: &domain.AccountCache{
				ID:     "acc-789",
				Status: domain.AccountStatusDeleted,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := domain.NewAccountCache(tt.id, tt.status)

			if cache.ID != tt.expected.ID {
				t.Errorf("Expected ID %s, got %s", tt.expected.ID, cache.ID)
			}
			if cache.Status != tt.expected.Status {
				t.Errorf("Expected Status %s, got %s", tt.expected.Status, cache.Status)
			}
		})
	}
}

func TestAccountCacheStatusMethods(t *testing.T) {
	t.Run("IsActive", func(t *testing.T) {
		tests := []struct {
			status   domain.AccountStatus
			expected bool
		}{
			{domain.AccountStatusActive, true},
			{domain.AccountStatusBlocked, false},
			{domain.AccountStatusDeleted, false},
		}

		for _, tt := range tests {
			cache := domain.NewAccountCache("acc-1", tt.status)
			if cache.IsActive() != tt.expected {
				t.Errorf("Status %s: expected IsActive()=%v, got %v", tt.status, tt.expected, cache.IsActive())
			}
		}
	})

	t.Run("IsDeleted", func(t *testing.T) {
		tests := []struct {
			status   domain.AccountStatus
			expected bool
		}{
			{domain.AccountStatusActive, false},
			{domain.AccountStatusBlocked, false},
			{domain.AccountStatusDeleted, true},
		}

		for _, tt := range tests {
			cache := domain.NewAccountCache("acc-1", tt.status)
			if cache.IsDeleted() != tt.expected {
				t.Errorf("Status %s: expected IsDeleted()=%v, got %v", tt.status, tt.expected, cache.IsDeleted())
			}
		}
	})

	t.Run("IsBlocked", func(t *testing.T) {
		tests := []struct {
			status   domain.AccountStatus
			expected bool
		}{
			{domain.AccountStatusActive, false},
			{domain.AccountStatusBlocked, true},
			{domain.AccountStatusDeleted, false},
		}

		for _, tt := range tests {
			cache := domain.NewAccountCache("acc-1", tt.status)
			if cache.IsBlocked() != tt.expected {
				t.Errorf("Status %s: expected IsBlocked()=%v, got %v", tt.status, tt.expected, cache.IsBlocked())
			}
		}
	})
}
