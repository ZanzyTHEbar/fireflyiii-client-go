package firefly

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/time/rate"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOAuth2Configuration(t *testing.T) {
	t.Run("OAuth2Config creation", func(t *testing.T) {
		config := &OAuth2Config{
			ClientID:     "test-client-id",
			ClientSecret: "test-client-secret",
			Scopes:       []string{"read", "write"},
			RedirectURL:  "http://localhost:8080/callback",
			AuthURL:      "https://example.com/auth",
			TokenURL:     "https://example.com/token",
		}

		assert.Equal(t, "test-client-id", config.ClientID)
		assert.Equal(t, "test-client-secret", config.ClientSecret)
		assert.Equal(t, []string{"read", "write"}, config.Scopes)
		assert.Equal(t, "http://localhost:8080/callback", config.RedirectURL)
		assert.Equal(t, "https://example.com/auth", config.AuthURL)
		assert.Equal(t, "https://example.com/token", config.TokenURL)
	})

	t.Run("ClientConfig with OAuth2", func(t *testing.T) {
		oauth2Config := OAuth2Config{
			ClientID:     "test-client",
			ClientSecret: "test-secret",
			TokenURL:     "https://example.com/token",
		}

		config := DefaultClientConfig().
			WithOAuth2(oauth2Config).
			WithTimeout(60*time.Second).
			WithRetry(5, 2*time.Second).
			WithRateLimit(100)

		assert.NotNil(t, config.OAuth2)
		assert.Equal(t, "test-client", config.OAuth2.ClientID)
		assert.Equal(t, 60*time.Second, config.Timeout)
		assert.Equal(t, 5, config.RetryCount)
		assert.Equal(t, 2*time.Second, config.RetryDelay)
		assert.Equal(t, 100, config.RateLimit)
	})
}

func TestRetryConfiguration(t *testing.T) {
	t.Run("DefaultRetryConfig", func(t *testing.T) {
		config := DefaultRetryConfig()

		assert.Equal(t, 3, config.MaxRetries)
		assert.Equal(t, time.Second, config.InitialDelay)
		assert.Equal(t, 30*time.Second, config.MaxDelay)
		assert.Equal(t, 2.0, config.BackoffFactor)
		assert.Contains(t, config.RetryableErrors, ErrNetwork)
		assert.Contains(t, config.RetryableErrors, ErrTimeout)
		assert.Contains(t, config.RetryableErrors, ErrServerError)
		assert.Contains(t, config.RetryableErrors, ErrRateLimit)
	})

	t.Run("isRetryableError", func(t *testing.T) {
		config := DefaultRetryConfig()

		// Test HTTP errors
		httpErr500 := &HTTPError{StatusCode: 500}
		assert.True(t, config.isRetryableError(httpErr500))

		httpErr429 := &HTTPError{StatusCode: 429}
		assert.True(t, config.isRetryableError(httpErr429))

		httpErr404 := &HTTPError{StatusCode: 404}
		assert.False(t, config.isRetryableError(httpErr404))

		// Test context errors
		assert.True(t, config.isRetryableError(context.DeadlineExceeded))
		assert.True(t, config.isRetryableError(context.Canceled))

		// Test nil error
		assert.False(t, config.isRetryableError(nil))
	})

	t.Run("calculateBackoffDelay", func(t *testing.T) {
		config := DefaultRetryConfig()

		// First attempt should return initial delay
		delay0 := config.calculateBackoffDelay(0)
		assert.Equal(t, time.Second, delay0)

		// Subsequent attempts should increase exponentially
		delay1 := config.calculateBackoffDelay(1)
		assert.True(t, delay1 >= 1800*time.Millisecond) // ~2s with jitter
		assert.True(t, delay1 <= 2200*time.Millisecond)

		delay2 := config.calculateBackoffDelay(2)
		assert.True(t, delay2 >= 3600*time.Millisecond) // ~4s with jitter
		assert.True(t, delay2 <= 4400*time.Millisecond)

		// Should not exceed max delay
		delay10 := config.calculateBackoffDelay(10)
		assert.True(t, delay10 <= 33*time.Second) // Max + 10% jitter
	})
}

func TestMiddleware(t *testing.T) {
	t.Run("MiddlewareChain", func(t *testing.T) {
		chain := NewMiddlewareChain()
		assert.NotNil(t, chain)
		assert.Empty(t, chain.middlewares)

		// Add logging middleware
		logger := func(format string, args ...interface{}) {
			t.Logf(format, args...)
		}
		loggingMW := NewLoggingMiddleware(logger)
		chain.Add(loggingMW)

		assert.Len(t, chain.middlewares, 1)
	})

	t.Run("LoggingMiddleware", func(t *testing.T) {
		var loggedMessages []string
		logger := func(format string, args ...interface{}) {
			loggedMessages = append(loggedMessages, format)
		}

		middleware := NewLoggingMiddleware(logger)
		ctx := context.Background()

		// Test request logging
		req := httptest.NewRequest("GET", "https://example.com/api/test", nil)
		processedReq, err := middleware.ProcessRequest(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, req, processedReq)
		assert.Contains(t, loggedMessages[0], "HTTP Request")

		// Test response logging
		resp := &http.Response{StatusCode: 200, Status: "200 OK"}
		processedResp, err := middleware.ProcessResponse(ctx, resp)
		require.NoError(t, err)
		assert.Equal(t, resp, processedResp)
		assert.Contains(t, loggedMessages[1], "HTTP Response")
	})

	t.Run("RateLimitMiddleware", func(t *testing.T) {
		// Create a very restrictive rate limiter for testing
		limiter := rate.NewLimiter(rate.Limit(1), 1) // 1 request per second, burst of 1
		middleware := NewRateLimitMiddleware(limiter)
		ctx := context.Background()

		req := httptest.NewRequest("GET", "https://example.com/api/test", nil)

		// First request should pass
		processedReq, err := middleware.ProcessRequest(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, req, processedReq)

		// Immediate second request should be rate limited (would block in real scenario)
		// For testing, we use a context with timeout
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
		defer cancel()

		_, err = middleware.ProcessRequest(ctxWithTimeout, req)
		assert.Error(t, err) // Should timeout due to rate limiting
	})

	t.Run("RetryMiddleware", func(t *testing.T) {
		config := DefaultRetryConfig()
		middleware := NewRetryMiddleware(config)
		ctx := context.Background()

		req := httptest.NewRequest("GET", "https://example.com/api/test", nil)

		// Test request pass-through
		processedReq, err := middleware.ProcessRequest(ctx, req)
		require.NoError(t, err)
		assert.Equal(t, req, processedReq)

		// Test successful response
		successResp := &http.Response{
			StatusCode: 200,
			Request:    req,
		}
		processedResp, err := middleware.ProcessResponse(ctx, successResp)
		require.NoError(t, err)
		assert.Equal(t, successResp, processedResp)

		// Test retryable error response
		retryableResp := &http.Response{
			StatusCode: 500,
			Request:    req,
		}
		_, err = middleware.ProcessResponse(ctx, retryableResp)
		assert.Error(t, err) // Should return error for retryable status
		assert.IsType(t, &HTTPError{}, err)

		// Test non-retryable error response
		nonRetryableResp := &http.Response{
			StatusCode: 404,
			Request:    req,
		}
		processedResp, err = middleware.ProcessResponse(ctx, nonRetryableResp)
		require.NoError(t, err)
		assert.Equal(t, nonRetryableResp, processedResp)
	})
}

func TestWebhookManager(t *testing.T) {
	t.Run("WebhookManager creation", func(t *testing.T) {
		manager := NewWebhookManager()
		assert.NotNil(t, manager)
		assert.NotNil(t, manager.handlers)
	})

	t.Run("RegisterHandler and ProcessWebhook", func(t *testing.T) {
		manager := NewWebhookManager()
		ctx := context.Background()

		var processedEvents []*WebhookEvent
		handler := WebhookHandlerFunc(func(ctx context.Context, event *WebhookEvent) error {
			processedEvents = append(processedEvents, event)
			return nil
		})

		// Register handler for transaction events
		manager.RegisterHandler("transaction.created", handler)

		// Test webhook processing
		payload := []byte(`{
			"id": "test-event-1",
			"type": "transaction.created",
			"timestamp": "2023-01-01T00:00:00Z",
			"data": {"transaction_id": "123"}
		}`)

		err := manager.ProcessWebhook(ctx, payload)
		require.NoError(t, err)
		assert.Len(t, processedEvents, 1)
		assert.Equal(t, "test-event-1", processedEvents[0].ID)
		assert.Equal(t, "transaction.created", processedEvents[0].Type)
	})

	t.Run("RegisterHandlerFunc", func(t *testing.T) {
		manager := NewWebhookManager()
		ctx := context.Background()

		var eventID string
		handlerFunc := func(ctx context.Context, event *WebhookEvent) error {
			eventID = event.ID
			return nil
		}

		manager.RegisterHandlerFunc("account.updated", handlerFunc)

		payload := []byte(`{
			"id": "test-event-2",
			"type": "account.updated",
			"timestamp": "2023-01-01T00:00:00Z",
			"data": {"account_id": "456"}
		}`)

		err := manager.ProcessWebhook(ctx, payload)
		require.NoError(t, err)
		assert.Equal(t, "test-event-2", eventID)
	})

	t.Run("No handlers registered", func(t *testing.T) {
		manager := NewWebhookManager()
		ctx := context.Background()

		payload := []byte(`{
			"id": "test-event-3",
			"type": "unknown.event",
			"timestamp": "2023-01-01T00:00:00Z",
			"data": {}
		}`)

		// Should not return error even if no handlers are registered
		err := manager.ProcessWebhook(ctx, payload)
		require.NoError(t, err)
	})

	t.Run("Invalid JSON payload", func(t *testing.T) {
		manager := NewWebhookManager()
		ctx := context.Background()

		payload := []byte(`{invalid json}`)

		err := manager.ProcessWebhook(ctx, payload)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to unmarshal webhook payload")
	})
}

func TestFireflyClientAdvancedFeatures(t *testing.T) {
	t.Run("NewFireflyClientWithConfig", func(t *testing.T) {
		config := DefaultClientConfig()
		config.BaseURL = "https://example.com"
		config.Token = "test-token"

		client, err := NewFireflyClientWithConfig(config)
		require.NoError(t, err)
		assert.NotNil(t, client)
		assert.Equal(t, "https://example.com", client.baseURL)
		assert.Equal(t, "test-token", client.token)
		assert.NotNil(t, client.middleware)
		assert.NotNil(t, client.webhookMgr)
	})

	t.Run("AddMiddleware", func(t *testing.T) {
		config := DefaultClientConfig()
		config.BaseURL = "https://example.com"
		config.Token = "test-token"

		client, err := NewFireflyClientWithConfig(config)
		require.NoError(t, err)

		logger := func(format string, args ...interface{}) {
			t.Logf(format, args...)
		}
		loggingMW := NewLoggingMiddleware(logger)

		client.AddMiddleware(loggingMW)
		assert.Len(t, client.middleware.middlewares, 1)
	})

	t.Run("GetWebhookManager", func(t *testing.T) {
		config := DefaultClientConfig()
		config.BaseURL = "https://example.com"
		config.Token = "test-token"

		client, err := NewFireflyClientWithConfig(config)
		require.NoError(t, err)

		manager := client.GetWebhookManager()
		assert.NotNil(t, manager)
		assert.Same(t, client.webhookMgr, manager)
	})

	t.Run("EnableDefaultMiddleware", func(t *testing.T) {
		config := DefaultClientConfig()
		config.BaseURL = "https://example.com"
		config.Token = "test-token"
		config.DebugMode = true
		config.RetryCount = 3

		client, err := NewFireflyClientWithConfig(config)
		require.NoError(t, err)

		initialCount := len(client.middleware.middlewares)
		client.EnableDefaultMiddleware()

		// Should have added rate limiting, logging, and retry middleware
		assert.Greater(t, len(client.middleware.middlewares), initialCount)
	})
}

func TestOAuth2Methods(t *testing.T) {
	t.Run("GenerateOAuth2AuthURL", func(t *testing.T) {
		config := DefaultClientConfig()
		config.BaseURL = "https://example.com"
		config.Token = "test-token"
		config.OAuth2 = &OAuth2Config{
			ClientID:    "test-client",
			AuthURL:     "https://example.com/auth",
			RedirectURL: "http://localhost:8080/callback",
		}

		client, err := NewFireflyClientWithConfig(config)
		require.NoError(t, err)

		// Test with provided state
		authURL, err := client.GenerateOAuth2AuthURL("test-state")
		require.NoError(t, err)
		assert.Contains(t, authURL, "https://example.com/auth")
		assert.Contains(t, authURL, "client_id=test-client")
		assert.Contains(t, authURL, "state=test-state")

		// Test with auto-generated state
		authURL2, err := client.GenerateOAuth2AuthURL("")
		require.NoError(t, err)
		assert.Contains(t, authURL2, "https://example.com/auth")
		assert.Contains(t, authURL2, "client_id=test-client")
		assert.Contains(t, authURL2, "state=")
	})

	t.Run("OAuth2 errors without configuration", func(t *testing.T) {
		config := DefaultClientConfig()
		config.BaseURL = "https://example.com"
		config.Token = "test-token"
		// No OAuth2 configuration

		client, err := NewFireflyClientWithConfig(config)
		require.NoError(t, err)

		ctx := context.Background()

		// Test client credentials without config
		_, err = client.GetOAuth2ClientCredentialsToken(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "OAuth2 Error")

		// Test auth URL without config
		_, err = client.GenerateOAuth2AuthURL("test-state")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "OAuth2 Error")

		// Test code exchange without config
		_, err = client.ExchangeOAuth2Code(ctx, "test-code", "test-state")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "OAuth2 Error")

		// Test token refresh without config
		_, err = client.RefreshOAuth2Token(ctx, "test-refresh-token")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "OAuth2 Error")
	})
}
