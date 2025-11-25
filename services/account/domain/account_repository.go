package domain

// AccountRepository defines the interface for account data operations
type AccountRepository interface {
	Create(account *Account) error
	GetByID(id string) (*Account, error)
	GetByAccountNumber(accountNumber string) (*Account, error)
	Update(account *Account) error
	Delete(id string) error
	List() ([]*Account, error)
}
