package application

// CreateAccountRequest represents the input data for creating an account
// Note: ID and AccountNumber are auto-generated using UUIDs
type CreateAccountRequest struct {
	BeholderName string `json:"beholder_name"`
	CountryCode  string `json:"country_code"`
}

// UpdateAccountRequest represents the input data for updating an account
type UpdateAccountRequest struct {
	ID            string `json:"id"`
	AccountNumber string `json:"account_number"`
	BeholderName  string `json:"beholder_name"`
	CountryCode   string `json:"country_code"`
	Status        string `json:"status"`
}

// AccountResponse represents the output data for account operations
type AccountResponse struct {
	ID            string `json:"id"`
	AccountNumber string `json:"account_number"`
	BeholderName  string `json:"beholder_name"`
	CountryCode   string `json:"country_code"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// AccountListResponse represents a list of accounts
type AccountListResponse struct {
	Accounts []AccountResponse `json:"accounts"`
	Total    int               `json:"total"`
}
