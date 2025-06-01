package firefly

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ZanzyTHEbar/errbuilder-go"
)

// Error codes specific to Firefly operations
const (
	ErrInvalidTransaction = "invalid_transaction"
	ErrInvalidAccount     = "invalid_account"
	ErrInvalidCategory    = "invalid_category"
	ErrInvalidBudget      = "invalid_budget"
	ErrInvalidAttachment  = "invalid_attachment"
	ErrAPIFailure         = "api_failure"
	ErrNotFound           = "not_found"
	ErrDuplicateEntry     = "duplicate_entry"
	ErrRateLimit          = "rate_limit"
	ErrAuthentication     = "authentication_failed"
	ErrAuthorization      = "authorization_failed"
	ErrNetwork            = "network_error"
	ErrTimeout            = "timeout_error"
	ErrServerError        = "server_error"
	ErrOAuth2             = "oauth2_error"
)

// HTTPError represents an HTTP-specific error with detailed context
type HTTPError struct {
	StatusCode   int               `json:"status_code"`
	Method       string            `json:"method"`
	URL          string            `json:"url"`
	Headers      map[string]string `json:"headers,omitempty"`
	Body         string            `json:"body,omitempty"`
	ResponseTime time.Duration     `json:"response_time"`
	Timestamp    time.Time         `json:"timestamp"`
}

// OAuth2Error represents OAuth2-specific errors
type OAuth2Error struct {
	ErrorCode        string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
	ErrorURI         string `json:"error_uri,omitempty"`
	State            string `json:"state,omitempty"`
}

// Error implements the error interface for HTTPError
func (h *HTTPError) Error() string {
	return fmt.Sprintf("HTTP %d: %s %s (took %v)", h.StatusCode, h.Method, h.URL, h.ResponseTime)
}

// Error implements the error interface for OAuth2Error
func (o *OAuth2Error) Error() string {
	if o.ErrorDescription != "" {
		return fmt.Sprintf("OAuth2 error: %s - %s", o.ErrorCode, o.ErrorDescription)
	}
	return fmt.Sprintf("OAuth2 error: %s", o.ErrorCode)
}

// NewHTTPError creates a new HTTP error with context
func NewHTTPError(statusCode int, method, url string, responseTime time.Duration) *HTTPError {
	return &HTTPError{
		StatusCode:   statusCode,
		Method:       method,
		URL:          url,
		ResponseTime: responseTime,
		Timestamp:    time.Now(),
	}
}

// WithHeaders adds headers to the HTTP error
func (h *HTTPError) WithHeaders(headers map[string]string) *HTTPError {
	h.Headers = headers
	return h
}

// WithBody adds response body to the HTTP error
func (h *HTTPError) WithBody(body string) *HTTPError {
	h.Body = body
	return h
}

// HTTPErrorFromResponse creates an HTTPError from an http.Response
func HTTPErrorFromResponse(resp *http.Response, method, url string, responseTime time.Duration) error {
	httpErr := NewHTTPError(resp.StatusCode, method, url, responseTime)

	// Add relevant headers
	headers := make(map[string]string)
	for _, key := range []string{"Content-Type", "X-Request-ID", "X-RateLimit-Remaining"} {
		if value := resp.Header.Get(key); value != "" {
			headers[key] = value
		}
	}
	if len(headers) > 0 {
		httpErr.WithHeaders(headers)
	}

	// Determine error type based on status code
	switch resp.StatusCode {
	case http.StatusUnauthorized:
		return AuthenticationErr(httpErr)
	case http.StatusForbidden:
		return AuthorizationErr(httpErr)
	case http.StatusNotFound:
		return NotFoundErr("Resource", httpErr)
	case http.StatusTooManyRequests:
		return RateLimitErr(httpErr)
	case http.StatusInternalServerError, http.StatusBadGateway, http.StatusServiceUnavailable:
		return ServerErr(httpErr)
	default:
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			return ClientErr(httpErr)
		}
		return ServerErr(httpErr)
	}
}

// AuthenticationErr returns an authentication error
func AuthenticationErr(err error) error {
	errs := make(errbuilder.ErrorMap)
	errs.Set("error_type", ErrAuthentication)
	errs.Set("help", "Check your API token or OAuth2 credentials")

	return errbuilder.NewErrBuilder().
		WithCode(errbuilder.CodeUnauthenticated).
		WithMsg("Authentication Failed").
		WithDetails(errbuilder.NewErrDetails(errs)).
		WithCause(err)
}

// AuthorizationErr returns an authorization error
func AuthorizationErr(err error) error {
	errs := make(errbuilder.ErrorMap)
	errs.Set("error_type", ErrAuthorization)
	errs.Set("help", "Check your permissions for this resource")

	return errbuilder.NewErrBuilder().
		WithCode(errbuilder.CodePermissionDenied).
		WithMsg("Authorization Failed").
		WithDetails(errbuilder.NewErrDetails(errs)).
		WithCause(err)
}

// NetworkErr returns a network error
func NetworkErr(err error) error {
	errs := make(errbuilder.ErrorMap)
	errs.Set("error_type", ErrNetwork)
	errs.Set("help", "Check your network connection and Firefly III URL")

	return errbuilder.NewErrBuilder().
		WithCode(errbuilder.CodeUnavailable).
		WithMsg("Network Error").
		WithDetails(errbuilder.NewErrDetails(errs)).
		WithCause(err)
}

// TimeoutErr returns a timeout error
func TimeoutErr(err error) error {
	errs := make(errbuilder.ErrorMap)
	errs.Set("error_type", ErrTimeout)
	errs.Set("help", "Request took too long to complete")

	return errbuilder.NewErrBuilder().
		WithCode(errbuilder.CodeDeadlineExceeded).
		WithMsg("Request Timeout").
		WithDetails(errbuilder.NewErrDetails(errs)).
		WithCause(err)
}

// ServerErr returns a server error
func ServerErr(err error) error {
	errs := make(errbuilder.ErrorMap)
	errs.Set("error_type", ErrServerError)
	errs.Set("help", "Firefly III server encountered an error")

	return errbuilder.NewErrBuilder().
		WithCode(errbuilder.CodeInternal).
		WithMsg("Server Error").
		WithDetails(errbuilder.NewErrDetails(errs)).
		WithCause(err)
}

// ClientErr returns a client error
func ClientErr(err error) error {
	errs := make(errbuilder.ErrorMap)
	errs.Set("error_type", ErrAPIFailure)
	errs.Set("help", "Check your request parameters")

	return errbuilder.NewErrBuilder().
		WithCode(errbuilder.CodeInvalidArgument).
		WithMsg("Client Error").
		WithDetails(errbuilder.NewErrDetails(errs)).
		WithCause(err)
}

// ContextErr returns a context-related error
func ContextErr(err error) error {
	if err == context.Canceled {
		return errbuilder.NewErrBuilder().
			WithCode(errbuilder.CodeCanceled).
			WithMsg("Request Cancelled").
			WithCause(err)
	}
	if err == context.DeadlineExceeded {
		return TimeoutErr(err)
	}
	return errbuilder.NewErrBuilder().
		WithCode(errbuilder.CodeInternal).
		WithMsg("Context Error").
		WithCause(err)
}

// OAuth2Err returns an OAuth2-specific error
func OAuth2Err(oauthErr *OAuth2Error) error {
	errs := make(errbuilder.ErrorMap)
	errs.Set("error_type", ErrOAuth2)
	errs.Set("oauth2_error", oauthErr.ErrorCode)
	errs.Set("oauth2_description", oauthErr.ErrorDescription)
	errs.Set("oauth2_uri", oauthErr.ErrorURI)

	return errbuilder.NewErrBuilder().
		WithCode(errbuilder.CodeUnauthenticated).
		WithMsg("OAuth2 Error").
		WithDetails(errbuilder.NewErrDetails(errs)).
		WithCause(oauthErr)
}

// TransactionValidationErr returns a validation error for transactions
func TransactionValidationErr(errors errbuilder.ErrorMap) error {
	return errbuilder.NewErrBuilder().
		WithCode(errbuilder.CodeInvalidArgument).
		WithMsg("Invalid Transaction Data").
		WithDetails(errbuilder.NewErrDetails(errors))
}

// AccountValidationErr returns a validation error for accounts
func AccountValidationErr(errors errbuilder.ErrorMap) error {
	return errbuilder.NewErrBuilder().
		WithCode(errbuilder.CodeInvalidArgument).
		WithMsg("Invalid Account Data").
		WithDetails(errbuilder.NewErrDetails(errors))
}

// CategoryValidationErr returns a validation error for categories
func CategoryValidationErr(errors errbuilder.ErrorMap) error {
	return errbuilder.NewErrBuilder().
		WithCode(errbuilder.CodeInvalidArgument).
		WithMsg("Invalid Category Data").
		WithDetails(errbuilder.NewErrDetails(errors))
}

// BudgetValidationErr returns a validation error for budgets
func BudgetValidationErr(errors errbuilder.ErrorMap) error {
	return errbuilder.NewErrBuilder().
		WithCode(errbuilder.CodeInvalidArgument).
		WithMsg("Invalid Budget Data").
		WithDetails(errbuilder.NewErrDetails(errors))
}

// AttachmentValidationErr returns a validation error for attachments
func AttachmentValidationErr(errors errbuilder.ErrorMap) error {
	return errbuilder.NewErrBuilder().
		WithCode(errbuilder.CodeInvalidArgument).
		WithMsg("Invalid Attachment Data").
		WithDetails(errbuilder.NewErrDetails(errors))
}

// APIErr returns an error for API failures
func APIErr(msg string, err error) error {
	return errbuilder.NewErrBuilder().
		WithCode(errbuilder.CodeInternal).
		WithMsg(msg).
		WithCause(err)
}

// NotFoundErr returns a not found error
func NotFoundErr(resourceType string, err error) error {
	return errbuilder.NewErrBuilder().
		WithCode(errbuilder.CodeNotFound).
		WithMsg(resourceType + " Not Found").
		WithCause(err)
}

// DuplicateErr returns a duplicate entry error
func DuplicateErr(resourceType string, err error) error {
	return errbuilder.NewErrBuilder().
		WithCode(errbuilder.CodeAlreadyExists).
		WithMsg("Duplicate " + resourceType).
		WithCause(err)
}

// RateLimitErr returns a rate limit error
func RateLimitErr(err error) error {
	return errbuilder.NewErrBuilder().
		WithCode(errbuilder.CodeResourceExhausted).
		WithMsg("Rate Limit Exceeded").
		WithCause(err)
}

// ValidationErr creates a generic validation error
func ValidationErr(entity string, errs errbuilder.ErrorMap) error {
	return errbuilder.NewErrBuilder().
		WithMsg(fmt.Sprintf("%s validation error: %v", entity, errs)).
		WithCode(errbuilder.CodeInvalidArgument).
		WithDetails(errbuilder.NewErrDetails(errs))
}

// validateTransaction validates a transaction and returns an error map
func validateTransaction(tx TransactionModel) errbuilder.ErrorMap {
	var errs errbuilder.ErrorMap

	if tx.Amount <= 0 {
		errs.Set("amount", "Amount must be greater than 0")
	}
	if tx.Currency == "" {
		errs.Set("currency", "Currency is required")
	}
	if tx.Description == "" {
		errs.Set("description", "Description is required")
	}
	if tx.TransType == "" {
		errs.Set("type", "Transaction type is required")
	}
	if tx.Date.IsZero() {
		errs.Set("date", "Date is required")
	}

	return errs
}

// validateAccount validates an account and returns an error map
func validateAccount(account AccountModel) errbuilder.ErrorMap {
	var errs errbuilder.ErrorMap

	if account.Name == "" {
		errs.Set("name", "Name is required")
	}
	if account.Type == "" {
		errs.Set("type", "Account type is required")
	}
	if account.Currency == "" {
		errs.Set("currency", "Currency is required")
	}

	return errs
}

// validateCategory validates a category and returns an error map
func validateCategory(category CategoryModel) errbuilder.ErrorMap {
	var errs errbuilder.ErrorMap

	if category.Name == "" {
		errs.Set("name", "Name is required")
	}

	return errs
}

// validateBudget validates a budget and returns an error map
func validateBudget(budget BudgetModel) errbuilder.ErrorMap {
	var errs errbuilder.ErrorMap

	if budget.Name == "" {
		errs.Set("name", "Name is required")
	}
	if budget.AutoBudgetAmount != nil && budget.AutoBudgetPeriod == nil {
		errs.Set("auto_budget_period", "Auto budget period is required when amount is set")
	}
	if budget.AutoBudgetPeriod != nil && budget.AutoBudgetAmount == nil {
		errs.Set("auto_budget_amount", "Auto budget amount is required when period is set")
	}

	return errs
}

// validateBudgetLimit validates a budget limit and returns an error map
func validateBudgetLimit(limit BudgetLimitModel) errbuilder.ErrorMap {
	var errs errbuilder.ErrorMap

	if limit.Amount == "" {
		errs.Set("amount", "Amount is required")
	}
	if limit.Period == "" {
		errs.Set("period", "Period is required")
	}
	if limit.Start.IsZero() {
		errs.Set("start", "Start date is required")
	}
	if limit.End.IsZero() {
		errs.Set("end", "End date is required")
	}
	if limit.End.Before(limit.Start) {
		errs.Set("end", "End date must be after start date")
	}

	return errs
}

// validateAttachment validates an attachment and returns an error map
func validateAttachment(filename string, file []byte, title string) errbuilder.ErrorMap {
	var errs errbuilder.ErrorMap

	if filename == "" {
		errs.Set("filename", "Filename is required")
	}
	if len(file) == 0 {
		errs.Set("file", "File content is required")
	}
	if title == "" {
		errs.Set("title", "Title is required")
	}

	return errs
}

// validatePiggyBank validates a piggy bank model
func validatePiggyBank(piggyBank PiggyBankModel) errbuilder.ErrorMap {
	var errs errbuilder.ErrorMap

	if piggyBank.Name == "" {
		errs.Set("name", "Name is required")
	}

	if piggyBank.TargetAmount == "" {
		errs.Set("target_amount", "Target amount is required")
	}

	if piggyBank.CurrencyCode == "" {
		errs.Set("currency_code", "Currency code is required")
	}

	if piggyBank.CurrencySymbol == "" {
		errs.Set("currency_symbol", "Currency symbol is required")
	}

	if piggyBank.TargetDate != nil && piggyBank.StartDate != nil && piggyBank.TargetDate.Before(*piggyBank.StartDate) {
		errs.Set("target_date", "Target date must be after start date")
	}

	return errs
}
