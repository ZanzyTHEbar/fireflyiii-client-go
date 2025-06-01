package firefly

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math"
	mathrand "math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"golang.org/x/time/rate"

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

// Middleware defines the interface for request/response middleware
type Middleware interface {
	ProcessRequest(ctx context.Context, req *http.Request) (*http.Request, error)
	ProcessResponse(ctx context.Context, resp *http.Response) (*http.Response, error)
}

// MiddlewareFunc allows using functions as middleware
type MiddlewareFunc func(ctx context.Context, req *http.Request, next func(*http.Request) (*http.Response, error)) (*http.Response, error)

// LoggingMiddleware logs HTTP requests and responses
type LoggingMiddleware struct {
	logger func(format string, args ...interface{})
}

// NewLoggingMiddleware creates a new logging middleware
func NewLoggingMiddleware(logger func(format string, args ...interface{})) *LoggingMiddleware {
	if logger == nil {
		logger = func(format string, args ...interface{}) {
			// Default no-op logger
		}
	}
	return &LoggingMiddleware{logger: logger}
}

// ProcessRequest logs the outgoing request
func (l *LoggingMiddleware) ProcessRequest(ctx context.Context, req *http.Request) (*http.Request, error) {
	l.logger("HTTP Request: %s %s", req.Method, req.URL.String())
	return req, nil
}

// ProcessResponse logs the incoming response
func (l *LoggingMiddleware) ProcessResponse(ctx context.Context, resp *http.Response) (*http.Response, error) {
	l.logger("HTTP Response: %d %s", resp.StatusCode, resp.Status)
	return resp, nil
}

// RetryMiddleware implements retry logic as middleware
type RetryMiddleware struct {
	config *RetryConfig
}

// NewRetryMiddleware creates a new retry middleware
func NewRetryMiddleware(config *RetryConfig) *RetryMiddleware {
	if config == nil {
		config = DefaultRetryConfig()
	}
	return &RetryMiddleware{config: config}
}

// ProcessRequest passes through the request
func (r *RetryMiddleware) ProcessRequest(ctx context.Context, req *http.Request) (*http.Request, error) {
	return req, nil
}

// ProcessResponse handles retry logic based on response
func (r *RetryMiddleware) ProcessResponse(ctx context.Context, resp *http.Response) (*http.Response, error) {
	// Create an HTTPError from the response for retry decision
	if resp.StatusCode >= 400 {
		httpErr := &HTTPError{
			StatusCode: resp.StatusCode,
			Method:     resp.Request.Method,
			URL:        resp.Request.URL.String(),
			Timestamp:  time.Now(),
		}

		// If this is a retryable error, return the error to trigger retry
		if r.config.isRetryableError(httpErr) {
			return nil, httpErr
		}
	}

	return resp, nil
}

// RateLimitMiddleware implements rate limiting as middleware
type RateLimitMiddleware struct {
	limiter *rate.Limiter
}

// NewRateLimitMiddleware creates a new rate limiting middleware
func NewRateLimitMiddleware(limiter *rate.Limiter) *RateLimitMiddleware {
	return &RateLimitMiddleware{limiter: limiter}
}

// ProcessRequest applies rate limiting before the request
func (r *RateLimitMiddleware) ProcessRequest(ctx context.Context, req *http.Request) (*http.Request, error) {
	if err := r.limiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rate limit wait failed: %w", err)
	}
	return req, nil
}

// ProcessResponse passes through the response
func (r *RateLimitMiddleware) ProcessResponse(ctx context.Context, resp *http.Response) (*http.Response, error) {
	return resp, nil
}

// MiddlewareChain manages a chain of middleware
type MiddlewareChain struct {
	middlewares []Middleware
	mu          sync.RWMutex
}

// NewMiddlewareChain creates a new middleware chain
func NewMiddlewareChain() *MiddlewareChain {
	return &MiddlewareChain{
		middlewares: make([]Middleware, 0),
	}
}

// Add adds middleware to the chain
func (m *MiddlewareChain) Add(middleware Middleware) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.middlewares = append(m.middlewares, middleware)
}

// ProcessRequest processes the request through all middleware
func (m *MiddlewareChain) ProcessRequest(ctx context.Context, req *http.Request) (*http.Request, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	currentReq := req
	for _, middleware := range m.middlewares {
		var err error
		currentReq, err = middleware.ProcessRequest(ctx, currentReq)
		if err != nil {
			return nil, err
		}
	}

	return currentReq, nil
}

// ProcessResponse processes the response through all middleware (in reverse order)
func (m *MiddlewareChain) ProcessResponse(ctx context.Context, resp *http.Response) (*http.Response, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	currentResp := resp
	// Process in reverse order
	for i := len(m.middlewares) - 1; i >= 0; i-- {
		var err error
		currentResp, err = m.middlewares[i].ProcessResponse(ctx, currentResp)
		if err != nil {
			return nil, err
		}
	}

	return currentResp, nil
}

// WebhookEvent represents a webhook event from Firefly III
type WebhookEvent struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Timestamp time.Time              `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// WebhookHandler defines the interface for handling webhook events
type WebhookHandler interface {
	HandleEvent(ctx context.Context, event *WebhookEvent) error
}

// WebhookHandlerFunc allows using functions as webhook handlers
type WebhookHandlerFunc func(ctx context.Context, event *WebhookEvent) error

// HandleEvent implements WebhookHandler for WebhookHandlerFunc
func (f WebhookHandlerFunc) HandleEvent(ctx context.Context, event *WebhookEvent) error {
	return f(ctx, event)
}

// WebhookManager manages webhook handlers and routing
type WebhookManager struct {
	handlers map[string][]WebhookHandler
	mu       sync.RWMutex
}

// NewWebhookManager creates a new webhook manager
func NewWebhookManager() *WebhookManager {
	return &WebhookManager{
		handlers: make(map[string][]WebhookHandler),
	}
}

// RegisterHandler registers a handler for a specific event type
func (w *WebhookManager) RegisterHandler(eventType string, handler WebhookHandler) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.handlers[eventType] = append(w.handlers[eventType], handler)
}

// RegisterHandlerFunc registers a handler function for a specific event type
func (w *WebhookManager) RegisterHandlerFunc(eventType string, handlerFunc func(ctx context.Context, event *WebhookEvent) error) {
	w.RegisterHandler(eventType, WebhookHandlerFunc(handlerFunc))
}

// ProcessWebhook processes an incoming webhook payload
func (w *WebhookManager) ProcessWebhook(ctx context.Context, payload []byte) error {
	var event WebhookEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		return fmt.Errorf("failed to unmarshal webhook payload: %w", err)
	}

	w.mu.RLock()
	handlers, exists := w.handlers[event.Type]
	w.mu.RUnlock()

	if !exists {
		// No handlers registered for this event type, not an error
		return nil
	}

	// Process handlers concurrently
	errChan := make(chan error, len(handlers))
	for _, handler := range handlers {
		go func(h WebhookHandler) {
			errChan <- h.HandleEvent(ctx, &event)
		}(handler)
	}

	// Collect errors
	var errs []error
	for i := 0; i < len(handlers); i++ {
		if err := <-errChan; err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("webhook processing errors: %v", errs)
	}

	return nil
}

// WebhookServer provides an HTTP server for receiving webhooks
type WebhookServer struct {
	manager *WebhookManager
	server  *http.Server
	secret  string
	path    string
}

// NewWebhookServer creates a new webhook server
func NewWebhookServer(addr, path, secret string, manager *WebhookManager) *WebhookServer {
	if manager == nil {
		manager = NewWebhookManager()
	}

	return &WebhookServer{
		manager: manager,
		secret:  secret,
		path:    path,
		server: &http.Server{
			Addr:         addr,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
	}
}

// Start starts the webhook server
func (ws *WebhookServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc(ws.path, ws.handleWebhook)
	ws.server.Handler = mux

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		ws.server.Shutdown(shutdownCtx)
	}()

	if err := ws.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("webhook server error: %w", err)
	}

	return nil
}

// handleWebhook handles incoming webhook requests
func (ws *WebhookServer) handleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	body := json.RawMessage(bodyBytes)

	// TODO: Implement webhook signature verification if secret is provided
	if ws.secret != "" {
		// Verify webhook signature here
		// This would typically involve checking HMAC signature in headers
	}

	ctx := r.Context()
	if err := ws.manager.ProcessWebhook(ctx, body); err != nil {
		http.Error(w, "Failed to process webhook", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// FireflyClient represents a client for the Firefly III API
type FireflyClient struct {
	baseURL    string
	token      string
	client     *http.Client
	clientAPI  *ClientWithResponses
	importers  map[string]importers.Importer
	config     *ClientConfig // Store configuration for advanced client
	middleware *MiddlewareChain
	webhookMgr *WebhookManager
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
	ID                  string                // The category ID
	Name                string                // The category name
	Notes               string                // Additional notes about the category
	CreatedAt           time.Time             // When the category was created
	UpdatedAt           time.Time             // When the category was last updated
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

// TODO: Add OAuth2 authentication configuration
type OAuth2Config struct {
	ClientID     string   `yaml:"client_id" json:"client_id"`
	ClientSecret string   `yaml:"client_secret" json:"client_secret"`
	Scopes       []string `yaml:"scopes" json:"scopes"`
	RedirectURL  string   `yaml:"redirect_url" json:"redirect_url"`
	AuthURL      string   `yaml:"auth_url" json:"auth_url"`
	TokenURL     string   `yaml:"token_url" json:"token_url"`
}

// ClientConfig holds configuration for the Firefly client
type ClientConfig struct {
	BaseURL    string        `yaml:"base_url" json:"base_url"`
	Token      string        `yaml:"token" json:"token"`
	Timeout    time.Duration `yaml:"timeout" json:"timeout"`
	RetryCount int           `yaml:"retry_count" json:"retry_count"`
	RetryDelay time.Duration `yaml:"retry_delay" json:"retry_delay"`
	RateLimit  int           `yaml:"rate_limit" json:"rate_limit"`
	OAuth2     *OAuth2Config `yaml:"oauth2,omitempty" json:"oauth2,omitempty"`
	UserAgent  string        `yaml:"user_agent" json:"user_agent"`
	DebugMode  bool          `yaml:"debug_mode" json:"debug_mode"`
}

// DefaultClientConfig returns a default client configuration
func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		Timeout:    30 * time.Second,
		RetryCount: 3,
		RetryDelay: time.Second,
		RateLimit:  60, // requests per minute
		UserAgent:  "firefly-client-go/1.0.0",
		DebugMode:  false,
	}
}

// WithOAuth2 configures OAuth2 authentication
func (c *ClientConfig) WithOAuth2(oauth2 OAuth2Config) *ClientConfig {
	c.OAuth2 = &oauth2
	return c
}

// WithTimeout sets the request timeout
func (c *ClientConfig) WithTimeout(timeout time.Duration) *ClientConfig {
	c.Timeout = timeout
	return c
}

// WithRetry configures retry behavior
func (c *ClientConfig) WithRetry(count int, delay time.Duration) *ClientConfig {
	c.RetryCount = count
	c.RetryDelay = delay
	return c
}

// WithRateLimit sets the rate limit (requests per minute)
func (c *ClientConfig) WithRateLimit(limit int) *ClientConfig {
	c.RateLimit = limit
	return c
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
		baseURL:    baseURL,
		token:      token,
		client:     client,
		clientAPI:  clientAPI,
		importers:  make(map[string]importers.Importer),
		middleware: NewMiddlewareChain(),
		webhookMgr: NewWebhookManager(),
	}, nil
}

// NewFireflyClientWithConfig creates a new Firefly III API client with advanced configuration
func NewFireflyClientWithConfig(config *ClientConfig) (*FireflyClient, error) {
	if config == nil {
		return nil, fmt.Errorf("client configuration cannot be nil")
	}

	if config.BaseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}

	// Create HTTP client with timeout and transport configuration
	client := &http.Client{
		Timeout: config.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: false,
		},
	}

	// Create request editor function for authentication and headers
	requestEditor := func(ctx context.Context, req *http.Request) error {
		// Add authentication
		if config.OAuth2 != nil {
			// TODO: Implement OAuth2 token refresh logic here
			// For now, fall back to token if available
			if config.Token != "" {
				req.Header.Set("Authorization", "Bearer "+config.Token)
			}
		} else if config.Token != "" {
			req.Header.Set("Authorization", "Bearer "+config.Token)
		}

		// Add user agent
		if config.UserAgent != "" {
			req.Header.Set("User-Agent", config.UserAgent)
		}

		// Add debug headers if enabled
		if config.DebugMode {
			req.Header.Set("X-Debug", "true")
		}

		return nil
	}

	// Create the generated client with responses and auth
	clientAPI, err := NewClientWithResponses(config.BaseURL, WithHTTPClient(client), WithRequestEditorFn(requestEditor))
	if err != nil {
		return nil, fmt.Errorf("failed to create Firefly III client: %w", err)
	}

	return &FireflyClient{
		baseURL:    config.BaseURL,
		token:      config.Token,
		client:     client,
		clientAPI:  clientAPI,
		importers:  make(map[string]importers.Importer),
		config:     config, // Store configuration for later use
		middleware: NewMiddlewareChain(),
		webhookMgr: NewWebhookManager(),
		limiter:    rate.NewLimiter(rate.Limit(config.RateLimit), 1), // requests per minute
	}, nil
}

// GetTransaction retrieves a single transaction by ID
func (c *FireflyClient) GetTransaction(ctx context.Context, id string) (*TransactionModel, error) {
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
func (c *FireflyClient) ListTransactions(ctx context.Context, page, limit int) ([]TransactionModel, error) {
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
func (c *FireflyClient) UpdateTransaction(ctx context.Context, id string, tx TransactionModel) error {
	// Validate transaction
	if errs := validateTransaction(tx); errs != nil {
		return TransactionValidationErr(errs)
	}

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
func (c *FireflyClient) DeleteTransaction(ctx context.Context, id string) error {
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
func (c *FireflyClient) SearchTransactions(ctx context.Context, query string) ([]TransactionModel, error) {
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
func (c *FireflyClient) ImportTransaction(ctx context.Context, tx TransactionModel) error {
	// Validate transaction
	if errs := validateTransaction(tx); errs != nil {
		return TransactionValidationErr(errs)
	}

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
func (c *FireflyClient) ImportTransactions(ctx context.Context, transactions []TransactionModel) error {
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
func (c *FireflyClient) CreateAccount(ctx context.Context, name, accountType, currency string) error {
	// Validate account
	account := AccountModel{
		Name:     name,
		Type:     accountType,
		Currency: currency,
	}
	if errs := validateAccount(account); errs != nil {
		return AccountValidationErr(errs)
	}

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
func (c *FireflyClient) UpdateBalance(ctx context.Context, accountID string, balance Balance) error {
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
func (c *FireflyClient) GetAccount(ctx context.Context, id string) (*AccountModel, error) {
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
func (c *FireflyClient) ListAccounts(ctx context.Context, page, limit int) ([]AccountModel, error) {
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
func (c *FireflyClient) DeleteAccount(ctx context.Context, id string) error {
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
func (c *FireflyClient) SearchAccounts(ctx context.Context, query string) ([]AccountModel, error) {
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
func (c *FireflyClient) CreateCategory(ctx context.Context, category CategoryModel) error {
	// Validate category
	if errs := validateCategory(category); errs != nil {
		return CategoryValidationErr(errs)
	}

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
func (c *FireflyClient) GetCategory(ctx context.Context, id string) (*CategoryModel, error) {
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
func (c *FireflyClient) ListCategories(ctx context.Context, page, limit int) ([]CategoryModel, error) {
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
func (c *FireflyClient) UpdateCategory(ctx context.Context, id string, category CategoryModel) error {
	// Validate category
	if errs := validateCategory(category); errs != nil {
		return CategoryValidationErr(errs)
	}

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
func (c *FireflyClient) DeleteCategory(ctx context.Context, id string) error {
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
func (c *FireflyClient) SearchCategories(ctx context.Context, query string) ([]CategoryModel, error) {
	// Get all categories (with a reasonable limit)
	categories, err := c.ListCategories(ctx, 1, 100)
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
func (c *FireflyClient) GetCategoryByName(ctx context.Context, name string) (*CategoryModel, error) {
	// Get all categories (with a reasonable limit)
	categories, err := c.ListCategories(ctx, 1, 100)
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

// NewFirefly creates a new Firefly III API client (convenience function)
// This is a convenience wrapper around NewFireflyClient for easier usage
func NewFirefly(baseURL, token string) *FireflyClient {
	client, err := NewFireflyClient(baseURL, token)
	if err != nil {
		// For backward compatibility, we'll panic on initialization errors
		// This matches the behavior expected by existing tests
		panic(fmt.Sprintf("failed to create Firefly client: %v", err))
	}
	return client
}

// OAuth2TokenResponse represents the response from OAuth2 token endpoint
type OAuth2TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// GetOAuth2ClientCredentialsToken obtains an access token using OAuth2 client credentials flow
func (c *FireflyClient) GetOAuth2ClientCredentialsToken(ctx context.Context) (*OAuth2TokenResponse, error) {
	if c.config == nil || c.config.OAuth2 == nil {
		return nil, OAuth2Err(&OAuth2Error{
			ErrorCode:        "oauth2_not_configured",
			ErrorDescription: "OAuth2 configuration is missing",
		})
	}

	oauth2Config := c.config.OAuth2
	if oauth2Config.ClientID == "" || oauth2Config.ClientSecret == "" || oauth2Config.TokenURL == "" {
		return nil, OAuth2Err(&OAuth2Error{
			ErrorCode:        "oauth2_configuration_incomplete",
			ErrorDescription: "client_id, client_secret, and token_url are required",
		})
	}

	// Use golang.org/x/oauth2/clientcredentials for client credentials flow
	config := &clientcredentials.Config{
		ClientID:     oauth2Config.ClientID,
		ClientSecret: oauth2Config.ClientSecret,
		TokenURL:     oauth2Config.TokenURL,
		Scopes:       oauth2Config.Scopes,
	}

	token, err := config.Token(ctx)
	if err != nil {
		return nil, OAuth2Err(&OAuth2Error{
			ErrorCode:        "token_request_failed",
			ErrorDescription: "Failed to obtain OAuth2 token: " + err.Error(),
		})
	}

	return &OAuth2TokenResponse{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		ExpiresIn:   int(time.Until(token.Expiry).Seconds()),
	}, nil
}

// GenerateOAuth2AuthURL generates an authorization URL for OAuth2 authorization code flow
func (c *FireflyClient) GenerateOAuth2AuthURL(state string) (string, error) {
	if c.config == nil || c.config.OAuth2 == nil {
		return "", OAuth2Err(&OAuth2Error{
			ErrorCode:        "oauth2_not_configured",
			ErrorDescription: "OAuth2 configuration is missing",
		})
	}

	oauth2Config := c.config.OAuth2
	if oauth2Config.ClientID == "" || oauth2Config.AuthURL == "" || oauth2Config.RedirectURL == "" {
		return "", OAuth2Err(&OAuth2Error{
			ErrorCode:        "oauth2_configuration_incomplete",
			ErrorDescription: "client_id, auth_url, and redirect_url are required",
		})
	}

	config := &oauth2.Config{
		ClientID:     oauth2Config.ClientID,
		ClientSecret: oauth2Config.ClientSecret,
		RedirectURL:  oauth2Config.RedirectURL,
		Scopes:       oauth2Config.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  oauth2Config.AuthURL,
			TokenURL: oauth2Config.TokenURL,
		},
	}

	// Generate a random state if not provided
	if state == "" {
		bytes := make([]byte, 32)
		if _, err := rand.Read(bytes); err != nil {
			return "", OAuth2Err(&OAuth2Error{
				ErrorCode:        "state_generation_failed",
				ErrorDescription: "Failed to generate state: " + err.Error(),
			})
		}
		state = base64.URLEncoding.EncodeToString(bytes)
	}

	return config.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

// ExchangeOAuth2Code exchanges an authorization code for access token
func (c *FireflyClient) ExchangeOAuth2Code(ctx context.Context, code, state string) (*OAuth2TokenResponse, error) {
	if c.config == nil || c.config.OAuth2 == nil {
		return nil, OAuth2Err(&OAuth2Error{
			ErrorCode:        "oauth2_not_configured",
			ErrorDescription: "OAuth2 configuration is missing",
		})
	}

	oauth2Config := c.config.OAuth2
	config := &oauth2.Config{
		ClientID:     oauth2Config.ClientID,
		ClientSecret: oauth2Config.ClientSecret,
		RedirectURL:  oauth2Config.RedirectURL,
		Scopes:       oauth2Config.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  oauth2Config.AuthURL,
			TokenURL: oauth2Config.TokenURL,
		},
	}

	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, OAuth2Err(&OAuth2Error{
			ErrorCode:        "code_exchange_failed",
			ErrorDescription: "Failed to exchange OAuth2 code: " + err.Error(),
		})
	}

	response := &OAuth2TokenResponse{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
		ExpiresIn:   int(time.Until(token.Expiry).Seconds()),
	}

	if token.RefreshToken != "" {
		response.RefreshToken = token.RefreshToken
	}

	return response, nil
}

// RefreshOAuth2Token refreshes an OAuth2 access token using refresh token
func (c *FireflyClient) RefreshOAuth2Token(ctx context.Context, refreshToken string) (*OAuth2TokenResponse, error) {
	if c.config == nil || c.config.OAuth2 == nil {
		return nil, OAuth2Err(&OAuth2Error{
			ErrorCode:        "oauth2_not_configured",
			ErrorDescription: "OAuth2 configuration is missing",
		})
	}

	oauth2Config := c.config.OAuth2
	config := &oauth2.Config{
		ClientID:     oauth2Config.ClientID,
		ClientSecret: oauth2Config.ClientSecret,
		RedirectURL:  oauth2Config.RedirectURL,
		Scopes:       oauth2Config.Scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  oauth2Config.AuthURL,
			TokenURL: oauth2Config.TokenURL,
		},
	}

	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}

	tokenSource := config.TokenSource(ctx, token)
	newToken, err := tokenSource.Token()
	if err != nil {
		return nil, OAuth2Err(&OAuth2Error{
			ErrorCode:        "token_refresh_failed",
			ErrorDescription: "Failed to refresh OAuth2 token: " + err.Error(),
		})
	}

	response := &OAuth2TokenResponse{
		AccessToken: newToken.AccessToken,
		TokenType:   newToken.TokenType,
		ExpiresIn:   int(time.Until(newToken.Expiry).Seconds()),
	}

	if newToken.RefreshToken != "" {
		response.RefreshToken = newToken.RefreshToken
	}

	return response, nil
}

// RetryConfig holds configuration for retry behavior
type RetryConfig struct {
	MaxRetries      int
	InitialDelay    time.Duration
	MaxDelay        time.Duration
	BackoffFactor   float64
	RetryableErrors []string
}

// DefaultRetryConfig returns a default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries:    3,
		InitialDelay:  time.Second,
		MaxDelay:      30 * time.Second,
		BackoffFactor: 2.0,
		RetryableErrors: []string{
			ErrNetwork,
			ErrTimeout,
			ErrServerError,
			ErrRateLimit,
		},
	}
}

// isRetryableError checks if an error is retryable based on configuration
func (r *RetryConfig) isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// Check for HTTP errors
	if httpErr, ok := err.(*HTTPError); ok {
		// Retry on 5xx server errors, 429 rate limit, and some 4xx errors
		switch httpErr.StatusCode {
		case http.StatusTooManyRequests, // 429
			http.StatusInternalServerError, // 500
			http.StatusBadGateway,          // 502
			http.StatusServiceUnavailable,  // 503
			http.StatusGatewayTimeout:      // 504
			return true
		case http.StatusRequestTimeout: // 408
			return true
		}
	}

	// Check for context errors (timeout, cancellation)
	if err == context.DeadlineExceeded || err == context.Canceled {
		return true
	}

	// Check error string for retryable error types
	errStr := err.Error()
	for _, retryableErr := range r.RetryableErrors {
		if strings.Contains(errStr, retryableErr) {
			return true
		}
	}

	return false
}

// calculateBackoffDelay calculates the delay for the next retry using exponential backoff
func (r *RetryConfig) calculateBackoffDelay(attempt int) time.Duration {
	if attempt <= 0 {
		return r.InitialDelay
	}

	delay := float64(r.InitialDelay) * math.Pow(r.BackoffFactor, float64(attempt))
	if delay > float64(r.MaxDelay) {
		delay = float64(r.MaxDelay)
	}

	// Add some jitter to avoid thundering herd
	jitter := delay * 0.1 * (0.5 - mathrand.Float64()) // ±10% jitter
	finalDelay := time.Duration(delay + jitter)

	return finalDelay
}

// RetryOperation wraps an operation with retry logic using exponential backoff
func (c *FireflyClient) RetryOperation(ctx context.Context, operation func(ctx context.Context) error) error {
	retryConfig := DefaultRetryConfig()
	if c.config != nil {
		retryConfig.MaxRetries = c.config.RetryCount
		retryConfig.InitialDelay = c.config.RetryDelay
	}

	var lastErr error
	for attempt := 0; attempt <= retryConfig.MaxRetries; attempt++ {
		// Check if context is done before attempting
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Execute the operation
		err := operation(ctx)
		if err == nil {
			return nil // Success
		}

		lastErr = err

		// Don't retry on the last attempt
		if attempt == retryConfig.MaxRetries {
			break
		}

		// Check if the error is retryable
		if !retryConfig.isRetryableError(err) {
			return err // Not retryable, return immediately
		}

		// Calculate backoff delay
		delay := retryConfig.calculateBackoffDelay(attempt)

		// Wait for the delay or context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
			// Continue to next attempt
		}
	}

	return lastErr
}

// AddMiddleware adds middleware to the client's middleware chain
func (c *FireflyClient) AddMiddleware(middleware Middleware) {
	c.middleware.Add(middleware)
}

// GetWebhookManager returns the client's webhook manager
func (c *FireflyClient) GetWebhookManager() *WebhookManager {
	return c.webhookMgr
}

// EnableDefaultMiddleware enables commonly used middleware with default configurations
func (c *FireflyClient) EnableDefaultMiddleware() {
	// Add rate limiting middleware
	if c.limiter != nil {
		c.AddMiddleware(NewRateLimitMiddleware(c.limiter))
	}

	// Add logging middleware if debug mode is enabled
	if c.config != nil && c.config.DebugMode {
		logger := func(format string, args ...interface{}) {
			fmt.Printf("[DEBUG] "+format+"\n", args...)
		}
		c.AddMiddleware(NewLoggingMiddleware(logger))
	}

	// Add retry middleware if retry is configured
	if c.config != nil && c.config.RetryCount > 0 {
		retryConfig := &RetryConfig{
			MaxRetries:    c.config.RetryCount,
			InitialDelay:  c.config.RetryDelay,
			MaxDelay:      30 * time.Second,
			BackoffFactor: 2.0,
			RetryableErrors: []string{
				ErrNetwork,
				ErrTimeout,
				ErrServerError,
				ErrRateLimit,
			},
		}
		c.AddMiddleware(NewRetryMiddleware(retryConfig))
	}
}
