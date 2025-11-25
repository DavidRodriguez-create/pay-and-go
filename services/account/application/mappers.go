package application

import (
	"time"

	"github.com/DavidRodriguez-create/pay-and-go/services/account/domain"
)

// ToAccountResponse converts a domain Account to an AccountResponse DTO
func ToAccountResponse(account *domain.Account) *AccountResponse {
	return &AccountResponse{
		ID:            account.ID,
		AccountNumber: account.AccountNumber,
		BeholderName:  account.BeholderName,
		CountryCode:   account.CountryCode,
		Status:        string(account.Status),
		CreatedAt:     account.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     account.UpdatedAt.Format(time.RFC3339),
	}
}

// ToAccountListResponse converts a slice of domain Accounts to an AccountListResponse DTO
func ToAccountListResponse(accounts []*domain.Account) *AccountListResponse {
	responses := make([]AccountResponse, len(accounts))
	for i, account := range accounts {
		responses[i] = *ToAccountResponse(account)
	}
	return &AccountListResponse{
		Accounts: responses,
		Total:    len(responses),
	}
}

// ToAccountStatus converts a string to domain.AccountStatus
func ToAccountStatus(status string) domain.AccountStatus {
	return domain.AccountStatus(status)
}
