package application

// GetAccountByID retrieves an account by its ID
func (s *AccountServiceImpl) GetAccountByID(id string) (*AccountResponse, error) {
	account, err := s.repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	return ToAccountResponse(account), nil
}

// GetAccountByAccountNumber retrieves an account by its account number
func (s *AccountServiceImpl) GetAccountByAccountNumber(accountNumber string) (*AccountResponse, error) {
	account, err := s.repository.GetByAccountNumber(accountNumber)
	if err != nil {
		return nil, err
	}
	return ToAccountResponse(account), nil
}

// ListAccounts retrieves all accounts
func (s *AccountServiceImpl) ListAccounts() (*AccountListResponse, error) {
	accounts, err := s.repository.List()
	if err != nil {
		return nil, err
	}
	return ToAccountListResponse(accounts), nil
}
