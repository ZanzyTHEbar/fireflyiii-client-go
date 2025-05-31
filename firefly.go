package firefly

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ZanzyTHEbar/fireflyiii-client-go/importers"
)

// TODO: Improve category operations to be more efficient by caching a dynamically generated/updated hashmap of categories as they are fetched

// FireflyClientInterface defines the interface for Firefly III API operations.
// This interface provides methods to interact with various resources in Firefly III
// such as transactions, accounts, categories, and more.
type FireflyClientInterface interface {
	// Transaction Operations

	// ImportTransaction creates a new transaction in Firefly III.
	// It takes a TransactionModel and returns an error if the operation fails.
	ImportTransaction(ctx context.Context, tx TransactionModel) error

	// ImportTransactions creates multiple transactions in Firefly III in a single operation.
	// It takes a slice of TransactionModel and returns an error if the operation fails.
	ImportTransactions(ctx context.Context, transactions []TransactionModel) error

	// GetTransaction retrieves a transaction by its ID.
	// It returns the transaction model and an error if the operation fails.
	GetTransaction(ctx context.Context, id string) (*TransactionModel, error)

	// ListTransactions retrieves a paginated list of transactions.
	// page: The page number to retrieve (starts at 1)
	// limit: The number of transactions per page
	// It returns a slice of transactions and an error if the operation fails.
	ListTransactions(ctx context.Context, page, limit int) ([]TransactionModel, error)

	// UpdateTransaction updates an existing transaction identified by id.
	// It takes the transaction ID and a TransactionModel with the updated values.
	// Returns an error if the operation fails.
	UpdateTransaction(ctx context.Context, id string, tx TransactionModel) error

	// DeleteTransaction removes a transaction from Firefly III.
	// It takes the transaction ID and returns an error if the operation fails.
	DeleteTransaction(ctx context.Context, id string) error

	// SearchTransactions searches for transactions matching the given query.
	// It returns a slice of matching transactions and an error if the operation fails.
	SearchTransactions(ctx context.Context, query string) ([]TransactionModel, error)

	// Account Operations

	// CreateAccount creates a new account in Firefly III.
	// name: The account name
	// accountType: The type of account (asset, expense, revenue, etc.)
	// currency: The currency code for the account
	// Returns an error if the operation fails.
	CreateAccount(ctx context.Context, name, accountType, currency string) error

	// UpdateBalance updates the balance of an account.
	// accountID: The ID of the account to update
	// balance: The new balance information
	// Returns an error if the operation fails.
	UpdateBalance(ctx context.Context, accountID string, balance Balance) error

	// GetAccount retrieves an account by its ID.
	// It returns the account model and an error if the operation fails.
	GetAccount(ctx context.Context, id string) (*AccountModel, error)

	// ListAccounts retrieves a paginated list of accounts.
	// page: The page number to retrieve (starts at 1)
	// limit: The number of accounts per page
	// It returns a slice of accounts and an error if the operation fails.
	ListAccounts(ctx context.Context, page, limit int) ([]AccountModel, error)

	// DeleteAccount removes an account from Firefly III.
	// It takes the account ID and returns an error if the operation fails.
	DeleteAccount(ctx context.Context, id string) error

	// SearchAccounts searches for accounts matching the given query.
	// It returns a slice of matching accounts and an error if the operation fails.
	SearchAccounts(ctx context.Context, query string) ([]AccountModel, error)

	// Category Operations

	// CreateCategory creates a new category in Firefly III.
	// It takes a CategoryModel and returns an error if the operation fails.
	CreateCategory(ctx context.Context, category CategoryModel) error

	// GetCategory retrieves a category by its ID.
	// It returns the category model and an error if the operation fails.
	GetCategory(ctx context.Context, id string) (*CategoryModel, error)

	// ListCategories retrieves a paginated list of categories.
	// page: The page number to retrieve (starts at 1)
	// limit: The number of categories per page
	// It returns a slice of categories and an error if the operation fails.
	ListCategories(ctx context.Context, page, limit int) ([]CategoryModel, error)

	// UpdateCategory updates an existing category identified by id.
	// It takes the category ID and a CategoryModel with the updated values.
	// Returns an error if the operation fails.
	UpdateCategory(ctx context.Context, id string, category CategoryModel) error

	// DeleteCategory removes a category from Firefly III.
	// It takes the category ID and returns an error if the operation fails.
	DeleteCategory(ctx context.Context, id string) error

	// SearchCategories searches for categories matching the given query.
	// It returns a slice of matching categories and an error if the operation fails.
	SearchCategories(ctx context.Context, query string) ([]CategoryModel, error)

	// GetCategoryByName retrieves a category by its name.
	// It returns the category model and an error if the operation fails.
	GetCategoryByName(ctx context.Context, name string) (*CategoryModel, error)

	// Attachment Operations

	// AddCategoryAttachment adds an attachment to a category.
	// categoryID: The ID of the category
	// filename: The name of the file
	// file: The file content as a byte slice
	// title: The title of the attachment
	// notes: Additional notes about the attachment
	// Returns the created attachment model and an error if the operation fails.
	AddCategoryAttachment(ctx context.Context, categoryID string, filename string, file []byte, title, notes string) (*AttachmentModel, error)

	// GetCategoryAttachments retrieves all attachments for a category.
	// It takes the category ID and returns a slice of attachments and an error if the operation fails.
	GetCategoryAttachments(ctx context.Context, categoryID string) ([]AttachmentModel, error)

	// DownloadCategoryAttachment downloads the content of an attachment.
	// It returns the file content as a byte slice, the filename, and an error if the operation fails.
	DownloadCategoryAttachment(ctx context.Context, attachmentID string) ([]byte, string, error)

	// DeleteCategoryAttachment removes an attachment from Firefly III.
	// It takes the attachment ID and returns an error if the operation fails.
	DeleteCategoryAttachment(ctx context.Context, attachmentID string) error

	// UpdateCategoryAttachment updates an existing attachment.
	// attachmentID: The ID of the attachment to update
	// filename: The new filename
	// title: The new title
	// notes: The new notes
	// Returns an error if the operation fails.
	UpdateCategoryAttachment(ctx context.Context, attachmentID string, filename, title, notes string) error

	// Budget Operations
	CreateBudget(budget BudgetModel) error
	GetBudget(id string) (*BudgetModel, error)
	ListBudgets(page, limit int) ([]BudgetModel, error)
	UpdateBudget(id string, budget BudgetModel) error
	DeleteBudget(id string) error
	SearchBudgets(query string) ([]BudgetModel, error)

	// Budget Limit Operations
	SetBudgetLimit(budgetID string, limit BudgetLimitModel) error
	GetBudgetLimits(budgetID string) ([]BudgetLimitModel, error)
	UpdateBudgetLimit(limitID string, limit BudgetLimitModel) error
	DeleteBudgetLimit(limitID string) error

	// Data Management Operations
	ExportData(dataType DataType, format ExportFormat) ([]byte, error)
	ImportData(dataType ImportType, format ImportFormat, data []byte, options *ImportOptions) (*ImportResult, error)
	DestroyData(dataType DataType) error
	BulkUpdateTransactions(query map[string]interface{}) error
	PurgeData() error

	// Importer Operations
	RegisterImporter(importer importers.Importer) error
	GetImporter(name string) (importers.Importer, error)
	ListImporters() []importers.Importer
	RunImporter(name string, options importers.ImportOptions) (*importers.ImportResult, error)
	GetImporterProgress(name string) (*importers.ImportProgress, error)
	CancelImporter(name string) error
}

// FireflyClient represents a client for the Firefly III API
type FireflyClient struct {
	baseURL    string
	token      string
	client     *http.Client
	clientAPI  *ClientWithResponses
	importers  map[string]importers.Importer
}

// TransactionModel represents a financial transaction in our domain model
type TransactionModel struct {
	ID              string
	Currency        string
	Amount          float64
	TransType       string // "deposit" or "withdrawal"
	Description     string
	Date            time.Time
	Category        string
	ForeignAmount   *float64
	ForeignCurrency *string
}

// AccountModel represents a financial account
type AccountModel struct {
	ID       string
	Name     string
	Type     string
	Currency string
	Balance  float64
	IBAN     string
	Number   string
	BankName string
	Active   bool
	Role     string
	Include  bool
}

// CategorySpentModel represents spending data for a category
type CategorySpentModel struct {
	Amount       string     `json:"amount"`
	CurrencyCode string     `json:"currency_code"`
	Date         *time.Time `json:"date"`
}

// CategoryEarnedModel represents earning data for a category
type CategoryEarnedModel struct {
	Amount       string     `json:"amount"`
	CurrencyCode string     `json:"currency_code"`
	Date         *time.Time `json:"date"`
}

// CategoryModel represents a category in Firefly III
type CategoryModel struct {
	ID                  string
	Name                string
	Notes               string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	Spent               []CategorySpentModel  // Total amount spent in this category
	Earned              []CategoryEarnedModel // Total amount earned in this category
	NativeCurrency      string                // The administration's native currency code
	NativeDecimalPlaces int32                 // The administration's native currency decimal places
	NativeSymbol        string                // The administration's native currency symbol
}

// Balance represents an account balance
type Balance struct {
	Currency string
	Amount   float64
}

// AttachmentModel represents a file attachment in our domain model
type AttachmentModel struct {
	ID          string
	Filename    string
	Title       string
	Notes       string
	Size        int32
	MimeType    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DownloadURL string
	Hash        string
}

// BudgetModel represents a budget in our domain model
type BudgetModel struct {
	ID               string
	Name             string
	Active           bool
	Notes            *string
	Order            *int32
	AutoBudgetAmount *string
	AutoBudgetPeriod *AutoBudgetPeriod
	AutoBudgetType   *AutoBudgetType
	Spent            *[]BudgetSpent
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// BudgetSpentModel represents spending within a budget period
type BudgetSpentModel struct {
	CurrencyCode string
	Amount       float64
	Period       string
}

// BudgetLimitModel represents a budget limit for a specific period
type BudgetLimitModel struct {
	ID        string
	BudgetID  *string
	Amount    string
	Period    string
	Start     time.Time
	End       time.Time
	Spent     *string
	Notes     *string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewFireflyClient creates a new Firefly III API client
func NewFireflyClient(baseURL, token string) (*FireflyClient, error) {
	// Create HTTP client with auth header
	client := &http.Client{}

	// Create the generated client with responses and auth
	clientAPI, err := NewClientWithResponses(baseURL, WithHTTPClient(client), WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	}))
	if err != nil {
		return nil, fmt.Errorf("failed to create Firefly III client: %w", err)
	}

	return &FireflyClient{
		baseURL:   baseURL,
		token:     token,
		client:    client,
		clientAPI: clientAPI,
		importers: make(map[string]importers.Importer),
	}, nil
}

// GetTransaction retrieves a single transaction by ID
func (c *FireflyClient) GetTransaction(id string) (*TransactionModel, error) {
	ctx := context.Background()

	// Call the API
	resp, err := c.clientAPI.GetTransactionWithResponse(ctx, id, &GetTransactionParams{})
	if err != nil {
		return nil, APIErr("Failed to get transaction", err)
	}

	// Check response
	if resp.StatusCode() == http.StatusNotFound {
		return nil, NotFoundErr("Transaction", fmt.Errorf("transaction not found: %s", id))
	}
	if resp.StatusCode() == http.StatusTooManyRequests {
		return nil, RateLimitErr(fmt.Errorf("rate limit exceeded"))
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, APIErr("Failed to get transaction", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	// Convert API response to TransactionModel
	if resp.HTTPResponse == nil || len(resp.Body) == 0 {
		return nil, APIErr("No transaction data found", fmt.Errorf("empty response"))
	}

	var apiResp TransactionSingle
	if err := json.Unmarshal(resp.Body, &apiResp); err != nil {
		return nil, APIErr("Failed to parse transaction response", err)
	}

	tx := &TransactionModel{
		ID:              apiResp.Data.Id,
		Description:     stringValue(apiResp.Data.Attributes.GroupTitle),
		Date:            *apiResp.Data.Attributes.CreatedAt,
		TransType:       apiResp.Data.Type,
		Category:        "",
		Currency:        "",
		Amount:          0,
		ForeignAmount:   nil,
		ForeignCurrency: nil,
	}

	// Handle amount and currency
	if len(apiResp.Data.Attributes.Transactions) > 0 {
		split := apiResp.Data.Attributes.Transactions[0]
		amount, err := strconv.ParseFloat(split.Amount, 64)
		if err != nil {
			return nil, APIErr("Failed to parse amount", err)
		}
		tx.Amount = amount
		if split.CurrencyCode != nil {
			tx.Currency = *split.CurrencyCode
		}

		// Handle foreign amount if present
		if split.ForeignAmount != nil {
			foreignAmount, err := strconv.ParseFloat(*split.ForeignAmount, 64)
			if err != nil {
				return nil, APIErr("Failed to parse foreign amount", err)
			}
			tx.ForeignAmount = float64Ptr(foreignAmount)
		}
		if split.ForeignCurrencyCode != nil {
			tx.ForeignCurrency = split.ForeignCurrencyCode
		}
	}

	return tx, nil
}

// ListTransactions retrieves a list of transactions with pagination
func (c *FireflyClient) ListTransactions(page, limit int) ([]TransactionModel, error) {
	ctx := context.Background()

	// Call the API
	resp, err := c.clientAPI.ListTransactionWithResponse(ctx, &ListTransactionParams{
		Page:  int32Ptr(page),
		Limit: int32Ptr(limit),
	})
	if err != nil {
		return nil, APIErr("Failed to list transactions", err)
	}

	// Check response
	if resp.StatusCode() == http.StatusTooManyRequests {
		return nil, RateLimitErr(fmt.Errorf("rate limit exceeded"))
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, APIErr("Failed to list transactions", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	// Convert API response to TransactionModels
	if resp.HTTPResponse == nil || len(resp.Body) == 0 {
		return []TransactionModel{}, nil
	}

	var apiResp TransactionArray
	if err := json.Unmarshal(resp.Body, &apiResp); err != nil {
		return nil, APIErr("Failed to parse transactions response", err)
	}

	transactions := make([]TransactionModel, 0, len(apiResp.Data))
	for _, txRead := range apiResp.Data {
		tx := TransactionModel{
			ID:              txRead.Id,
			Description:     stringValue(txRead.Attributes.GroupTitle),
			Date:            *txRead.Attributes.CreatedAt,
			TransType:       txRead.Type,
			Category:        "",
			Currency:        "",
			Amount:          0,
			ForeignAmount:   nil,
			ForeignCurrency: nil,
		}

		// Handle amount and currency
		if len(txRead.Attributes.Transactions) > 0 {
			split := txRead.Attributes.Transactions[0]
			amount, err := strconv.ParseFloat(split.Amount, 64)
			if err != nil {
				return nil, APIErr("Failed to parse amount", err)
			}
			tx.Amount = amount
			if split.CurrencyCode != nil {
				tx.Currency = *split.CurrencyCode
			}

			// Handle foreign amount if present
			if split.ForeignAmount != nil {
				foreignAmount, err := strconv.ParseFloat(*split.ForeignAmount, 64)
				if err != nil {
					return nil, APIErr("Failed to parse foreign amount", err)
				}
				tx.ForeignAmount = float64Ptr(foreignAmount)
			}
			if split.ForeignCurrencyCode != nil {
				tx.ForeignCurrency = split.ForeignCurrencyCode
			}
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}

// UpdateTransaction updates an existing transaction
func (c *FireflyClient) UpdateTransaction(id string, tx TransactionModel) error {
	// Validate transaction
	if errs := validateTransaction(tx); errs != nil {
		return TransactionValidationErr(errs)
	}

	ctx := context.Background()
	txType := TransactionTypeProperty(tx.TransType)

	// Convert our transaction to the API format
	apiTx := UpdateTransactionJSONRequestBody{
		ApplyRules: boolPtr(true),
		Transactions: &[]TransactionSplitUpdate{{
			Type:         &txType,
			Date:         timePtr(tx.Date),
			Amount:       stringPtr(fmt.Sprintf("%.2f", tx.Amount)),
			Description:  stringPtr(tx.Description),
			CurrencyCode: stringPtr(tx.Currency),
			CategoryName: &tx.Category,
		}},
	}

	// Handle foreign amount if present
	if tx.ForeignAmount != nil && tx.ForeignCurrency != nil {
		(*apiTx.Transactions)[0].ForeignAmount = stringPtr(fmt.Sprintf("%.2f", *tx.ForeignAmount))
		(*apiTx.Transactions)[0].ForeignCurrencyCode = tx.ForeignCurrency
	}

	// Call the API
	resp, err := c.clientAPI.UpdateTransactionWithResponse(ctx, id, &UpdateTransactionParams{}, apiTx)
	if err != nil {
		return APIErr("Failed to update transaction", err)
	}

	// Check response
	if resp.StatusCode() == http.StatusNotFound {
		return NotFoundErr("Transaction", fmt.Errorf("transaction not found: %s", id))
	}
	if resp.StatusCode() == http.StatusTooManyRequests {
		return RateLimitErr(fmt.Errorf("rate limit exceeded"))
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return APIErr("Failed to update transaction", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	return nil
}

// DeleteTransaction deletes a transaction by ID
func (c *FireflyClient) DeleteTransaction(id string) error {
	ctx := context.Background()

	// Call the API
	resp, err := c.clientAPI.DeleteTransactionWithResponse(ctx, id, &DeleteTransactionParams{})
	if err != nil {
		return APIErr("Failed to delete transaction", err)
	}

	// Check response
	if resp.StatusCode() == http.StatusNotFound {
		return NotFoundErr("Transaction", fmt.Errorf("transaction not found: %s", id))
	}
	if resp.StatusCode() == http.StatusTooManyRequests {
		return RateLimitErr(fmt.Errorf("rate limit exceeded"))
	}
	if resp.StatusCode() != http.StatusNoContent {
		return APIErr("Failed to delete transaction", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	return nil
}

// SearchTransactions searches for transactions matching the query
func (c *FireflyClient) SearchTransactions(query string) ([]TransactionModel, error) {
	ctx := context.Background()

	// Call the API
	resp, err := c.clientAPI.SearchTransactionsWithResponse(ctx, &SearchTransactionsParams{
		Query: query,
	})
	if err != nil {
		return nil, APIErr("Failed to search transactions", err)
	}

	// Check response
	if resp.StatusCode() == http.StatusTooManyRequests {
		return nil, RateLimitErr(fmt.Errorf("rate limit exceeded"))
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, APIErr("Failed to search transactions", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	// Convert API response to TransactionModels
	if resp.HTTPResponse == nil || len(resp.Body) == 0 {
		return []TransactionModel{}, nil
	}

	var apiResp TransactionArray
	if err := json.Unmarshal(resp.Body, &apiResp); err != nil {
		return nil, APIErr("Failed to parse transactions response", err)
	}

	transactions := make([]TransactionModel, 0, len(apiResp.Data))
	for _, txRead := range apiResp.Data {
		tx := TransactionModel{
			ID:              txRead.Id,
			Description:     stringValue(txRead.Attributes.GroupTitle),
			Date:            *txRead.Attributes.CreatedAt,
			TransType:       txRead.Type,
			Category:        "",
			Currency:        "",
			Amount:          0,
			ForeignAmount:   nil,
			ForeignCurrency: nil,
		}

		// Handle amount and currency
		if len(txRead.Attributes.Transactions) > 0 {
			split := txRead.Attributes.Transactions[0]
			amount, err := strconv.ParseFloat(split.Amount, 64)
			if err != nil {
				return nil, APIErr("Failed to parse amount", err)
			}
			tx.Amount = amount
			if split.CurrencyCode != nil {
				tx.Currency = *split.CurrencyCode
			}

			// Handle foreign amount if present
			if split.ForeignAmount != nil {
				foreignAmount, err := strconv.ParseFloat(*split.ForeignAmount, 64)
				if err != nil {
					return nil, APIErr("Failed to parse foreign amount", err)
				}
				tx.ForeignAmount = float64Ptr(foreignAmount)
			}
			if split.ForeignCurrencyCode != nil {
				tx.ForeignCurrency = split.ForeignCurrencyCode
			}
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}

// ImportTransaction imports a single transaction
func (c *FireflyClient) ImportTransaction(tx TransactionModel) error {
	// Validate transaction
	if errs := validateTransaction(tx); errs != nil {
		return TransactionValidationErr(errs)
	}

	ctx := context.Background()
	txType := TransactionTypeProperty(tx.TransType)

	// Convert our transaction to the API format
	apiTx := StoreTransactionJSONRequestBody{
		ErrorIfDuplicateHash: boolPtr(true),
		ApplyRules:           boolPtr(true),
		Transactions: []TransactionSplitStore{
			{
				Type:         txType,
				Date:         tx.Date,
				Amount:       fmt.Sprintf("%.2f", tx.Amount),
				Description:  tx.Description,
				CurrencyCode: stringPtr(tx.Currency),
				CategoryName: &tx.Category,
			},
		},
	}

	// Handle foreign amount if present
	if tx.ForeignAmount != nil && tx.ForeignCurrency != nil {
		apiTx.Transactions[0].ForeignAmount = stringPtr(fmt.Sprintf("%.2f", *tx.ForeignAmount))
		apiTx.Transactions[0].ForeignCurrencyCode = tx.ForeignCurrency
	}

	// Call the API
	resp, err := c.clientAPI.StoreTransactionWithResponse(ctx, &StoreTransactionParams{}, apiTx)
	if err != nil {
		return APIErr("Failed to import transaction", err)
	}

	// Check response
	if resp.StatusCode() == http.StatusConflict {
		return DuplicateErr("Transaction", fmt.Errorf("transaction already exists"))
	}
	if resp.StatusCode() == http.StatusTooManyRequests {
		return RateLimitErr(fmt.Errorf("rate limit exceeded"))
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return APIErr("Failed to import transaction", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	return nil
}

// ImportTransactions imports multiple transactions in batch
func (c *FireflyClient) ImportTransactions(transactions []TransactionModel) error {
	ctx := context.Background()

	// Validate all transactions first
	for _, tx := range transactions {
		if errs := validateTransaction(tx); errs != nil {
			return TransactionValidationErr(errs)
		}
	}

	// Convert transactions to API format
	splits := make([]TransactionSplitStore, len(transactions))
	for i, tx := range transactions {
		txType := TransactionTypeProperty(tx.TransType)
		splits[i] = TransactionSplitStore{
			Type:         txType,
			Date:         tx.Date,
			Amount:       fmt.Sprintf("%.2f", tx.Amount),
			Description:  tx.Description,
			CurrencyCode: stringPtr(tx.Currency),
			CategoryName: &tx.Category,
		}

		// Handle foreign amount if present
		if tx.ForeignAmount != nil && tx.ForeignCurrency != nil {
			splits[i].ForeignAmount = stringPtr(fmt.Sprintf("%.2f", *tx.ForeignAmount))
			splits[i].ForeignCurrencyCode = tx.ForeignCurrency
		}
	}

	// Create batch request
	apiTx := StoreTransactionJSONRequestBody{
		ErrorIfDuplicateHash: boolPtr(true),
		ApplyRules:           boolPtr(true),
		Transactions:         splits,
	}

	// Call the API
	resp, err := c.clientAPI.StoreTransactionWithResponse(ctx, &StoreTransactionParams{}, apiTx)
	if err != nil {
		return APIErr("Failed to import transactions", err)
	}

	// Check response
	if resp.StatusCode() == http.StatusConflict {
		return DuplicateErr("Transaction", fmt.Errorf("one or more transactions already exist"))
	}
	if resp.StatusCode() == http.StatusTooManyRequests {
		return RateLimitErr(fmt.Errorf("rate limit exceeded"))
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return APIErr("Failed to import transactions", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	return nil
}

// CreateAccount creates a new account
func (c *FireflyClient) CreateAccount(name, accountType, currency string) error {
	// Validate account
	account := AccountModel{
		Name:     name,
		Type:     accountType,
		Currency: currency,
	}
	if errs := validateAccount(account); errs != nil {
		return AccountValidationErr(errs)
	}

	ctx := context.Background()

	// Create account request
	accountRequest := StoreAccountJSONRequestBody{
		Name:         name,
		Type:         ShortAccountTypeProperty(accountType),
		CurrencyCode: stringPtr(currency),
	}

	// Call the API
	resp, err := c.clientAPI.StoreAccountWithResponse(ctx, &StoreAccountParams{}, accountRequest)
	if err != nil {
		return APIErr("Failed to create account", err)
	}

	// Check response
	if resp.StatusCode() == http.StatusConflict {
		return DuplicateErr("Account", fmt.Errorf("account already exists"))
	}
	if resp.StatusCode() == http.StatusTooManyRequests {
		return RateLimitErr(fmt.Errorf("rate limit exceeded"))
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return APIErr("Failed to create account", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	return nil
}

// UpdateBalance updates an account's balance
func (c *FireflyClient) UpdateBalance(accountID string, balance Balance) error {
	ctx := context.Background()

	// Convert float64 to string for API
	balanceStr := fmt.Sprintf("%.2f", balance.Amount)

	// Create balance update request
	update := UpdateAccountJSONRequestBody{
		CurrencyCode:   stringPtr(balance.Currency),
		OpeningBalance: stringPtr(balanceStr),
	}

	// Call the API
	resp, err := c.clientAPI.UpdateAccountWithResponse(ctx, accountID, &UpdateAccountParams{}, update)
	if err != nil {
		return APIErr("Failed to update balance", err)
	}

	// Check response
	if resp.StatusCode() == http.StatusNotFound {
		return NotFoundErr("Account", fmt.Errorf("account not found: %s", accountID))
	}
	if resp.StatusCode() == http.StatusTooManyRequests {
		return RateLimitErr(fmt.Errorf("rate limit exceeded"))
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return APIErr("Failed to update balance", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	return nil
}

// GetAccount retrieves a single account by ID
func (c *FireflyClient) GetAccount(id string) (*AccountModel, error) {
	ctx := context.Background()

	// Call the API
	resp, err := c.clientAPI.GetAccountWithResponse(ctx, id, &GetAccountParams{})
	if err != nil {
		return nil, APIErr("Failed to get account", err)
	}

	// Check response
	if resp.StatusCode() == http.StatusNotFound {
		return nil, NotFoundErr("Account", fmt.Errorf("account not found: %s", id))
	}
	if resp.StatusCode() == http.StatusTooManyRequests {
		return nil, RateLimitErr(fmt.Errorf("rate limit exceeded"))
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, APIErr("Failed to get account", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	// Convert API response to AccountModel
	if resp.HTTPResponse == nil || len(resp.Body) == 0 {
		return nil, APIErr("No account data found", fmt.Errorf("empty response"))
	}

	var apiResp AccountSingle
	if err := json.Unmarshal(resp.Body, &apiResp); err != nil {
		return nil, APIErr("Failed to parse account response", err)
	}

	// Parse balance
	balance := float64(0)
	if apiResp.Data.Attributes.CurrentBalance != nil {
		var err error
		balance, err = strconv.ParseFloat(*apiResp.Data.Attributes.CurrentBalance, 64)
		if err != nil {
			return nil, APIErr("Failed to parse balance", err)
		}
	}

	// Get account role
	role := ""
	if apiResp.Data.Attributes.AccountRole != nil {
		role = string(*apiResp.Data.Attributes.AccountRole)
	}

	account := &AccountModel{
		ID:       apiResp.Data.Id,
		Name:     apiResp.Data.Attributes.Name,
		Type:     string(apiResp.Data.Attributes.Type),
		Currency: stringValue(apiResp.Data.Attributes.CurrencyCode),
		Balance:  balance,
		IBAN:     stringValue(apiResp.Data.Attributes.Iban),
		Number:   stringValue(apiResp.Data.Attributes.AccountNumber),
		BankName: "", // Not available in API
		Active:   boolValue(apiResp.Data.Attributes.Active),
		Role:     role,
		Include:  boolValue(apiResp.Data.Attributes.IncludeNetWorth),
	}

	return account, nil
}

// ListAccounts retrieves a list of accounts with pagination
func (c *FireflyClient) ListAccounts(page, limit int) ([]AccountModel, error) {
	ctx := context.Background()

	// Call the API
	resp, err := c.clientAPI.ListAccountWithResponse(ctx, &ListAccountParams{
		Page:  int32Ptr(page),
		Limit: int32Ptr(limit),
	})
	if err != nil {
		return nil, APIErr("Failed to list accounts", err)
	}

	// Check response
	if resp.StatusCode() == http.StatusTooManyRequests {
		return nil, RateLimitErr(fmt.Errorf("rate limit exceeded"))
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, APIErr("Failed to list accounts", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	// Convert API response to AccountModels
	if resp.HTTPResponse == nil || len(resp.Body) == 0 {
		return []AccountModel{}, nil
	}

	var apiResp AccountArray
	if err := json.Unmarshal(resp.Body, &apiResp); err != nil {
		return nil, APIErr("Failed to parse accounts response", err)
	}

	accounts := make([]AccountModel, 0, len(apiResp.Data))
	for _, accountRead := range apiResp.Data {
		// Parse balance
		balance := float64(0)
		if accountRead.Attributes.CurrentBalance != nil {
			var err error
			balance, err = strconv.ParseFloat(*accountRead.Attributes.CurrentBalance, 64)
			if err != nil {
				return nil, APIErr("Failed to parse balance", err)
			}
		}

		// Get account role
		role := ""
		if accountRead.Attributes.AccountRole != nil {
			role = string(*accountRead.Attributes.AccountRole)
		}

		account := AccountModel{
			ID:       accountRead.Id,
			Name:     accountRead.Attributes.Name,
			Type:     string(accountRead.Attributes.Type),
			Currency: stringValue(accountRead.Attributes.CurrencyCode),
			Balance:  balance,
			IBAN:     stringValue(accountRead.Attributes.Iban),
			Number:   stringValue(accountRead.Attributes.AccountNumber),
			BankName: "", // Not available in API
			Active:   boolValue(accountRead.Attributes.Active),
			Role:     role,
			Include:  boolValue(accountRead.Attributes.IncludeNetWorth),
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

// DeleteAccount deletes an account by ID
func (c *FireflyClient) DeleteAccount(id string) error {
	ctx := context.Background()

	// Call the API
	resp, err := c.clientAPI.DeleteAccountWithResponse(ctx, id, &DeleteAccountParams{})
	if err != nil {
		return APIErr("Failed to delete account", err)
	}

	// Check response
	if resp.StatusCode() != http.StatusNoContent {
		return APIErr("Failed to delete account", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	return nil
}

// SearchAccounts searches for accounts matching the query
func (c *FireflyClient) SearchAccounts(query string) ([]AccountModel, error) {
	ctx := context.Background()

	// Call the API
	resp, err := c.clientAPI.SearchAccountsWithResponse(ctx, &SearchAccountsParams{
		Query: query,
		Field: AccountSearchFieldFilter("all"), // Search in all fields
	})
	if err != nil {
		return nil, APIErr("Failed to search accounts", err)
	}

	// Check response
	if resp.StatusCode() == http.StatusTooManyRequests {
		return nil, RateLimitErr(fmt.Errorf("rate limit exceeded"))
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, APIErr("Failed to search accounts", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	// Convert API response to AccountModels
	if resp.HTTPResponse == nil || len(resp.Body) == 0 {
		return []AccountModel{}, nil
	}

	var apiResp AccountArray
	if err := json.Unmarshal(resp.Body, &apiResp); err != nil {
		return nil, APIErr("Failed to parse accounts response", err)
	}

	accounts := make([]AccountModel, 0, len(apiResp.Data))
	for _, accountRead := range apiResp.Data {
		// Parse balance
		balance := float64(0)
		if accountRead.Attributes.CurrentBalance != nil {
			var err error
			balance, err = strconv.ParseFloat(*accountRead.Attributes.CurrentBalance, 64)
			if err != nil {
				return nil, APIErr("Failed to parse balance", err)
			}
		}

		// Get account role
		role := ""
		if accountRead.Attributes.AccountRole != nil {
			role = string(*accountRead.Attributes.AccountRole)
		}

		account := AccountModel{
			ID:       accountRead.Id,
			Name:     accountRead.Attributes.Name,
			Type:     string(accountRead.Attributes.Type),
			Currency: stringValue(accountRead.Attributes.CurrencyCode),
			Balance:  balance,
			IBAN:     stringValue(accountRead.Attributes.Iban),
			Number:   stringValue(accountRead.Attributes.AccountNumber),
			BankName: "", // Not available in API
			Active:   boolValue(accountRead.Attributes.Active),
			Role:     role,
			Include:  boolValue(accountRead.Attributes.IncludeNetWorth),
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

// CreateCategory creates a new category
func (c *FireflyClient) CreateCategory(category CategoryModel) error {
	// Validate category
	if errs := validateCategory(category); errs != nil {
		return CategoryValidationErr(errs)
	}

	ctx := context.Background()

	notes := category.Notes // Create a copy to get address of
	// Create category request
	categoryRequest := StoreCategoryJSONRequestBody{
		Name:  category.Name,
		Notes: &notes,
	}

	// Call the API
	resp, err := c.clientAPI.StoreCategoryWithResponse(ctx, &StoreCategoryParams{}, categoryRequest)
	if err != nil {
		return APIErr("Failed to create category", err)
	}

	// Check response
	if resp.StatusCode() == http.StatusConflict {
		return DuplicateErr("Category", fmt.Errorf("category already exists"))
	}
	if resp.StatusCode() == http.StatusTooManyRequests {
		return RateLimitErr(fmt.Errorf("rate limit exceeded"))
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return APIErr("Failed to create category", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	return nil
}

// GetCategory retrieves a single category by ID
func (c *FireflyClient) GetCategory(id string) (*CategoryModel, error) {
	ctx := context.Background()
	response, err := c.clientAPI.GetCategoryWithResponse(ctx, id, &GetCategoryParams{})
	if err != nil {
		return nil, APIErr("Failed to get category", err)
	}

	if response.StatusCode() == http.StatusNotFound {
		return nil, NotFoundErr("Category", fmt.Errorf("category not found: %s", id))
	}
	if response.StatusCode() == http.StatusTooManyRequests {
		return nil, RateLimitErr(fmt.Errorf("rate limit exceeded"))
	}
	if response.StatusCode() != http.StatusOK {
		return nil, APIErr("Failed to get category", fmt.Errorf("unexpected status: %s", response.Status()))
	}

	if response.HTTPResponse == nil || len(response.Body) == 0 {
		return nil, APIErr("No category data found", fmt.Errorf("empty response"))
	}

	var apiResp CategorySingle
	if err := json.Unmarshal(response.Body, &apiResp); err != nil {
		return nil, APIErr("Failed to parse category response", err)
	}

	category := &CategoryModel{
		ID:             apiResp.Data.Id,
		Name:           apiResp.Data.Attributes.Name,
		Notes:          stringValue(apiResp.Data.Attributes.Notes),
		Spent:          make([]CategorySpentModel, 0),
		Earned:         make([]CategoryEarnedModel, 0),
		CreatedAt:      timeValue(apiResp.Data.Attributes.CreatedAt),
		UpdatedAt:      timeValue(apiResp.Data.Attributes.UpdatedAt),
		NativeCurrency: stringValue(apiResp.Data.Attributes.NativeCurrencyCode),
		NativeSymbol:   stringValue(apiResp.Data.Attributes.NativeCurrencySymbol),
	}

	// Process spent amounts
	if apiResp.Data.Attributes.Spent != nil {
		for _, spent := range *apiResp.Data.Attributes.Spent {
			category.Spent = append(category.Spent, CategorySpentModel{
				Amount:       stringValue(spent.Sum),
				CurrencyCode: stringValue(spent.CurrencyCode),
				Date:         nil, // API doesn't provide transaction date in category response
			})
		}
	}

	// Process earned amounts
	if apiResp.Data.Attributes.Earned != nil {
		for _, earned := range *apiResp.Data.Attributes.Earned {
			category.Earned = append(category.Earned, CategoryEarnedModel{
				Amount:       stringValue(earned.Sum),
				CurrencyCode: stringValue(earned.CurrencyCode),
				Date:         nil, // API doesn't provide transaction date in category response
			})
		}
	}

	return category, nil
}

// ListCategories retrieves a list of categories with pagination
func (c *FireflyClient) ListCategories(page, limit int) ([]CategoryModel, error) {
	ctx := context.Background()

	// Convert page and limit to int32
	page32 := int32(page)
	limit32 := int32(limit)

	// Call the API
	resp, err := c.clientAPI.ListCategoryWithResponse(ctx, &ListCategoryParams{
		Page:  &page32,
		Limit: &limit32,
	})
	if err != nil {
		return nil, APIErr("Failed to list categories", err)
	}

	// Check response
	if resp.StatusCode() == http.StatusTooManyRequests {
		return nil, RateLimitErr(fmt.Errorf("rate limit exceeded"))
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, APIErr("Failed to list categories", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	// Convert API response to CategoryModel array
	if resp.HTTPResponse == nil || len(resp.Body) == 0 {
		return nil, APIErr("No category data found", fmt.Errorf("empty response"))
	}

	var apiResp CategoryArray
	if err := json.Unmarshal(resp.Body, &apiResp); err != nil {
		return nil, APIErr("Failed to parse categories response", err)
	}

	categories := make([]CategoryModel, 0, len(apiResp.Data))
	for _, categoryRead := range apiResp.Data {
		category := CategoryModel{
			ID:                  categoryRead.Id,
			Name:                categoryRead.Attributes.Name,
			Notes:               stringValue(categoryRead.Attributes.Notes),
			CreatedAt:           timeValue(categoryRead.Attributes.CreatedAt),
			UpdatedAt:           timeValue(categoryRead.Attributes.UpdatedAt),
			Spent:               []CategorySpentModel{},
			Earned:              []CategoryEarnedModel{},
			NativeCurrency:      stringValue(categoryRead.Attributes.NativeCurrencyCode),
			NativeDecimalPlaces: int32Value(categoryRead.Attributes.NativeCurrencyDecimalPlaces),
			NativeSymbol:        stringValue(categoryRead.Attributes.NativeCurrencySymbol),
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// UpdateCategory updates an existing category
func (c *FireflyClient) UpdateCategory(id string, category CategoryModel) error {
	// Validate category
	if errs := validateCategory(category); errs != nil {
		return CategoryValidationErr(errs)
	}

	ctx := context.Background()

	notes := category.Notes // Create a copy to get address of
	update := UpdateCategoryJSONRequestBody{
		Name:  category.Name,
		Notes: &notes,
	}

	// Call the API
	resp, err := c.clientAPI.UpdateCategoryWithResponse(ctx, id, &UpdateCategoryParams{}, update)
	if err != nil {
		return APIErr("Failed to update category", err)
	}

	// Check response
	if resp.StatusCode() == http.StatusNotFound {
		return NotFoundErr("Category", fmt.Errorf("category not found: %s", id))
	}
	if resp.StatusCode() == http.StatusConflict {
		return DuplicateErr("Category", fmt.Errorf("category with this name already exists"))
	}
	if resp.StatusCode() == http.StatusTooManyRequests {
		return RateLimitErr(fmt.Errorf("rate limit exceeded"))
	}
	if resp.StatusCode() != http.StatusOK {
		return APIErr("Failed to update category", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	return nil
}

// DeleteCategory deletes a category
func (c *FireflyClient) DeleteCategory(id string) error {
	ctx := context.Background()

	// Call the API
	resp, err := c.clientAPI.DeleteCategoryWithResponse(ctx, id, &DeleteCategoryParams{})
	if err != nil {
		return APIErr("Failed to delete category", err)
	}

	// Check response
	if resp.StatusCode() == http.StatusNotFound {
		return NotFoundErr("Category", fmt.Errorf("category not found: %s", id))
	}
	if resp.StatusCode() == http.StatusTooManyRequests {
		return RateLimitErr(fmt.Errorf("rate limit exceeded"))
	}
	if resp.StatusCode() != http.StatusNoContent {
		return APIErr("Failed to delete category", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	return nil
}

// SearchCategories searches for categories matching the query
func (c *FireflyClient) SearchCategories(query string) ([]CategoryModel, error) {
	// Get all categories (with a reasonable limit)
	categories, err := c.ListCategories(1, 100)
	if err != nil {
		return nil, APIErr("Failed to search categories", err)
	}

	// Filter categories based on the query (case-insensitive)
	query = strings.ToLower(query)
	var results []CategoryModel
	for _, category := range categories {
		if strings.Contains(strings.ToLower(category.Name), query) ||
			strings.Contains(strings.ToLower(category.Notes), query) {
			results = append(results, category)
		}
	}

	return results, nil
}

// GetCategoryByName retrieves a category by its exact name (case-insensitive)
func (c *FireflyClient) GetCategoryByName(name string) (*CategoryModel, error) {
	// Get all categories (with a reasonable limit)
	categories, err := c.ListCategories(1, 100)
	if err != nil {
		return nil, APIErr("Failed to get category by name", err)
	}

	// Find the category with matching name (case-insensitive)
	name = strings.ToLower(name)
	for _, category := range categories {
		if strings.ToLower(category.Name) == name {
			return &category, nil
		}
	}

	return nil, NotFoundErr("Category", fmt.Errorf("category not found: %s", name))
}

// CreateBudget creates a new budget
func (c *FireflyClient) CreateBudget(budget BudgetModel) error {
	// Validate budget
	if errs := validateBudget(budget); errs != nil {
		return BudgetValidationErr(errs)
	}

	ctx := context.Background()

	// Create budget request
	budgetRequest := StoreBudgetJSONRequestBody{
		Name:             budget.Name,
		Active:           boolPtr(budget.Active),
		Notes:            budget.Notes,
		Order:            budget.Order,
		AutoBudgetAmount: budget.AutoBudgetAmount,
		AutoBudgetPeriod: budget.AutoBudgetPeriod,
		AutoBudgetType:   budget.AutoBudgetType,
	}

	// Call the API
	resp, err := c.clientAPI.StoreBudgetWithResponse(ctx, &StoreBudgetParams{}, budgetRequest)
	if err != nil {
		return APIErr("Failed to create budget", err)
	}

	// Check response
	if resp.StatusCode() == http.StatusConflict {
		return DuplicateErr("Budget", fmt.Errorf("budget already exists"))
	}
	if resp.StatusCode() == http.StatusTooManyRequests {
		return RateLimitErr(fmt.Errorf("rate limit exceeded"))
	}
	if resp.StatusCode() != http.StatusOK && resp.StatusCode() != http.StatusCreated {
		return APIErr("Failed to create budget", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	return nil
}

// GetBudget retrieves a single budget by ID
func (c *FireflyClient) GetBudget(id string) (*BudgetModel, error) {
	ctx := context.Background()

	// Call the API
	resp, err := c.clientAPI.GetBudgetWithResponse(ctx, id, &GetBudgetParams{})
	if err != nil {
		return nil, APIErr("Failed to get budget", err)
	}

	// Check response
	if resp.StatusCode() == http.StatusNotFound {
		return nil, NotFoundErr("Budget", fmt.Errorf("budget not found: %s", id))
	}
	if resp.StatusCode() == http.StatusTooManyRequests {
		return nil, RateLimitErr(fmt.Errorf("rate limit exceeded"))
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, APIErr("Failed to get budget", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	// Convert API response to BudgetModel
	if resp.HTTPResponse == nil || len(resp.Body) == 0 {
		return nil, APIErr("No budget data found", fmt.Errorf("empty response"))
	}

	var apiResp BudgetSingle
	if err := json.Unmarshal(resp.Body, &apiResp); err != nil {
		return nil, APIErr("Failed to parse budget response", err)
	}

	budget := &BudgetModel{
		ID:               apiResp.Data.Id,
		Name:             apiResp.Data.Attributes.Name,
		Active:           boolValue(apiResp.Data.Attributes.Active),
		Notes:            apiResp.Data.Attributes.Notes,
		Order:            apiResp.Data.Attributes.Order,
		AutoBudgetAmount: apiResp.Data.Attributes.AutoBudgetAmount,
		AutoBudgetPeriod: apiResp.Data.Attributes.AutoBudgetPeriod,
		AutoBudgetType:   apiResp.Data.Attributes.AutoBudgetType,
		Spent:            apiResp.Data.Attributes.Spent,
		CreatedAt:        timeValue(apiResp.Data.Attributes.CreatedAt),
		UpdatedAt:        timeValue(apiResp.Data.Attributes.UpdatedAt),
	}

	return budget, nil
}

// ListBudgets retrieves a list of budgets with pagination
func (c *FireflyClient) ListBudgets(page, limit int) ([]BudgetModel, error) {
	ctx := context.Background()

	// Call the API
	resp, err := c.clientAPI.ListBudgetWithResponse(ctx, &ListBudgetParams{
		Page:  int32Ptr(page),
		Limit: int32Ptr(limit),
	})
	if err != nil {
		return nil, APIErr("Failed to list budgets", err)
	}

	// Check response
	if resp.StatusCode() == http.StatusTooManyRequests {
		return nil, RateLimitErr(fmt.Errorf("rate limit exceeded"))
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, APIErr("Failed to list budgets", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	// Convert API response to BudgetModel array
	if resp.HTTPResponse == nil || len(resp.Body) == 0 {
		return nil, APIErr("No budget data found", fmt.Errorf("empty response"))
	}

	var apiResp BudgetArray
	if err := json.Unmarshal(resp.Body, &apiResp); err != nil {
		return nil, APIErr("Failed to parse budgets response", err)
	}

	budgets := make([]BudgetModel, 0, len(apiResp.Data))
	for _, budgetRead := range apiResp.Data {
		budget := BudgetModel{
			ID:               budgetRead.Id,
			Name:             budgetRead.Attributes.Name,
			Active:           boolValue(budgetRead.Attributes.Active),
			Notes:            budgetRead.Attributes.Notes,
			Order:            budgetRead.Attributes.Order,
			AutoBudgetAmount: budgetRead.Attributes.AutoBudgetAmount,
			AutoBudgetPeriod: budgetRead.Attributes.AutoBudgetPeriod,
			AutoBudgetType:   budgetRead.Attributes.AutoBudgetType,
			Spent:            budgetRead.Attributes.Spent,
			CreatedAt:        timeValue(budgetRead.Attributes.CreatedAt),
			UpdatedAt:        timeValue(budgetRead.Attributes.UpdatedAt),
		}
		budgets = append(budgets, budget)
	}

	return budgets, nil
}

// UpdateBudget updates an existing budget
func (c *FireflyClient) UpdateBudget(id string, budget BudgetModel) error {
	// Validate budget
	if errs := validateBudget(budget); errs != nil {
		return BudgetValidationErr(errs)
	}

	ctx := context.Background()

	// Create budget update request
	update := UpdateBudgetJSONRequestBody{
		Name:             budget.Name,
		Active:           boolPtr(budget.Active),
		Notes:            budget.Notes,
		Order:            budget.Order,
		AutoBudgetAmount: budget.AutoBudgetAmount,
		AutoBudgetPeriod: budget.AutoBudgetPeriod,
		AutoBudgetType:   budget.AutoBudgetType,
	}

	// Call the API
	resp, err := c.clientAPI.UpdateBudgetWithResponse(ctx, id, &UpdateBudgetParams{}, update)
	if err != nil {
		return APIErr("Failed to update budget", err)
	}

	// Check response
	if resp.StatusCode() == http.StatusNotFound {
		return NotFoundErr("Budget", fmt.Errorf("budget not found: %s", id))
	}
	if resp.StatusCode() == http.StatusTooManyRequests {
		return RateLimitErr(fmt.Errorf("rate limit exceeded"))
	}
	if resp.StatusCode() != http.StatusOK {
		return APIErr("Failed to update budget", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	return nil
}

// DeleteBudget deletes a budget
func (c *FireflyClient) DeleteBudget(id string) error {
	ctx := context.Background()

	// Call the API
	resp, err := c.clientAPI.DeleteBudgetWithResponse(ctx, id, &DeleteBudgetParams{})
	if err != nil {
		return APIErr("Failed to delete budget", err)
	}

	// Check response
	if resp.StatusCode() == http.StatusNotFound {
		return NotFoundErr("Budget", fmt.Errorf("budget not found: %s", id))
	}
	if resp.StatusCode() == http.StatusTooManyRequests {
		return RateLimitErr(fmt.Errorf("rate limit exceeded"))
	}
	if resp.StatusCode() != http.StatusNoContent {
		return APIErr("Failed to delete budget", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	return nil
}

// SetBudgetLimit sets a budget limit
func (c *FireflyClient) SetBudgetLimit(budgetID string, limit BudgetLimitModel) error {
	// Validate budget limit
	if errs := validateBudgetLimit(limit); errs != nil {
		return BudgetValidationErr(errs)
	}

	ctx := context.Background()

	// Create budget limit update request
	update := UpdateBudgetLimitJSONRequestBody{
		Amount: limit.Amount,
		Period: stringPtr(limit.Period),
		Start:  limit.Start,
		End:    limit.End,
		Notes:  limit.Notes,
	}

	// Call the API
	resp, err := c.clientAPI.UpdateBudgetLimitWithResponse(ctx, budgetID, limit.ID, &UpdateBudgetLimitParams{}, update)
	if err != nil {
		return APIErr("Failed to update budget limit", err)
	}

	// Check response
	if resp.StatusCode() == http.StatusNotFound {
		return NotFoundErr("Budget Limit", fmt.Errorf("budget limit not found: %s", limit.ID))
	}
	if resp.StatusCode() == http.StatusTooManyRequests {
		return RateLimitErr(fmt.Errorf("rate limit exceeded"))
	}
	if resp.StatusCode() != http.StatusOK {
		return APIErr("Failed to update budget limit", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	return nil
}

// GetBudgetLimits retrieves all budget limits for a budget
func (c *FireflyClient) GetBudgetLimits(budgetID string) ([]BudgetLimitModel, error) {
	ctx := context.Background()

	// Call the API
	resp, err := c.clientAPI.ListBudgetLimitWithResponse(ctx, &ListBudgetLimitParams{})
	if err != nil {
		return nil, APIErr("Failed to list budget limits", err)
	}

	// Check response
	if resp.StatusCode() == http.StatusNotFound {
		return nil, NotFoundErr("Budget", fmt.Errorf("budget not found: %s", budgetID))
	}
	if resp.StatusCode() == http.StatusTooManyRequests {
		return nil, RateLimitErr(fmt.Errorf("rate limit exceeded"))
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, APIErr("Failed to list budget limits", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	// Convert API response to BudgetLimitModel array
	if resp.HTTPResponse == nil || len(resp.Body) == 0 {
		return nil, APIErr("No budget limit data found", fmt.Errorf("empty response"))
	}

	var apiResp BudgetLimitArray
	if err := json.Unmarshal(resp.Body, &apiResp); err != nil {
		return nil, APIErr("Failed to parse budget limits response", err)
	}

	limits := make([]BudgetLimitModel, 0, len(apiResp.Data))
	for _, limitRead := range apiResp.Data {
		limit := BudgetLimitModel{
			ID:        limitRead.Id,
			BudgetID:  limitRead.Attributes.BudgetId,
			Amount:    limitRead.Attributes.Amount,
			Period:    stringValue(limitRead.Attributes.Period),
			Start:     limitRead.Attributes.Start,
			End:       limitRead.Attributes.End,
			Spent:     limitRead.Attributes.Spent,
			Notes:     limitRead.Attributes.Notes,
			CreatedAt: timeValue(limitRead.Attributes.CreatedAt),
			UpdatedAt: timeValue(limitRead.Attributes.UpdatedAt),
		}
		limits = append(limits, limit)
	}

	return limits, nil
}

// UpdateBudgetLimit updates an existing budget limit
func (c *FireflyClient) UpdateBudgetLimit(limitID string, limit BudgetLimitModel) error {
	// Validate budget limit
	if errs := validateBudgetLimit(limit); errs != nil {
		return BudgetValidationErr(errs)
	}

	ctx := context.Background()

	// Create budget limit update request
	update := UpdateBudgetLimitJSONRequestBody{
		Amount: limit.Amount,
		Period: stringPtr(limit.Period),
		Start:  limit.Start,
		End:    limit.End,
		Notes:  limit.Notes,
	}

	// Call the API
	resp, err := c.clientAPI.UpdateBudgetLimitWithResponse(ctx, stringValue(limit.BudgetID), limitID, &UpdateBudgetLimitParams{}, update)
	if err != nil {
		return APIErr("Failed to update budget limit", err)
	}

	// Check response
	if resp.StatusCode() == http.StatusNotFound {
		return NotFoundErr("Budget Limit", fmt.Errorf("budget limit not found: %s", limitID))
	}
	if resp.StatusCode() == http.StatusTooManyRequests {
		return RateLimitErr(fmt.Errorf("rate limit exceeded"))
	}
	if resp.StatusCode() != http.StatusOK {
		return APIErr("Failed to update budget limit", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	return nil
}

// DeleteBudgetLimit deletes a budget limit
func (c *FireflyClient) DeleteBudgetLimit(limitID string) error {
	ctx := context.Background()

	// Get the budget limit first to get its budget ID
	limits, err := c.GetBudgetLimits("")
	if err != nil {
		return fmt.Errorf("failed to get budget limit info: %w", err)
	}

	// Find the budget ID for this limit
	var budgetID string
	for _, limit := range limits {
		if limit.ID == limitID && limit.BudgetID != nil {
			budgetID = *limit.BudgetID
			break
		}
	}

	if budgetID == "" {
		return fmt.Errorf("could not find budget ID for limit: %s", limitID)
	}

	// Call the API
	resp, err := c.clientAPI.DeleteBudgetLimitWithResponse(ctx, budgetID, limitID, &DeleteBudgetLimitParams{})
	if err != nil {
		return APIErr("Failed to delete budget limit", err)
	}

	// Check response
	switch resp.StatusCode() {
	case http.StatusNotFound:
		return NotFoundErr("Budget Limit", fmt.Errorf("budget limit not found: %s", limitID))
	case http.StatusTooManyRequests:
		return RateLimitErr(fmt.Errorf("rate limit exceeded"))
	case http.StatusNoContent:
		// Successful response, continue
	default:
		return APIErr("Failed to delete budget limit", fmt.Errorf("unexpected status: %s", resp.Status()))
	}

	return nil
}

// RegisterImporter registers a new importer
func (c *FireflyClient) RegisterImporter(importer importers.Importer) error {
	config := importers.ImporterConfig{}
	if err := importer.ValidateConfig(config); err != nil {
		return fmt.Errorf("invalid importer configuration: %w", err)
	}

	c.importers[config.Name] = importer
	return nil
}

// GetImporter retrieves a registered importer by name
func (c *FireflyClient) GetImporter(name string) (importers.Importer, error) {
	importer, exists := c.importers[name]
	if !exists {
		return nil, fmt.Errorf("importer not found: %s", name)
	}
	return importer, nil
}

// ListImporters returns all registered importers
func (c *FireflyClient) ListImporters() []importers.Importer {
	importerList := make([]importers.Importer, 0, len(c.importers))
	for _, importer := range c.importers {
		importerList = append(importerList, importer)
	}
	return importerList
}

// RunImporter runs an importer with the given options
func (c *FireflyClient) RunImporter(name string, options importers.ImportOptions) (*importers.ImportResult, error) {
	importer, err := c.GetImporter(name)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	return importer.Import(ctx, options)
}

// GetImporterProgress gets the progress of an importer
func (c *FireflyClient) GetImporterProgress(name string) (*importers.ImportProgress, error) {
	importer, err := c.GetImporter(name)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	return importer.GetProgress(ctx)
}

// CancelImporter cancels an importer's operation
func (c *FireflyClient) CancelImporter(name string) error {
	importer, err := c.GetImporter(name)
	if err != nil {
		return err
	}

	ctx := context.Background()
	return importer.Cancel(ctx)
}
