package firefly

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestNewClient tests the creation of a new client
func TestNewClient(t *testing.T) {
	baseURL := "https://example.com/api"
	token := "test-token"
	
	client := NewClient(baseURL, token)
	
	if client == nil {
		t.Fatal("Expected client to be non-nil")
	}
	
	// Test client properties
	// Add assertions based on the actual implementation
}

// mockServer creates a test HTTP server that returns the given status code and response body
func mockServer(t *testing.T, status int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, err := w.Write([]byte(body))
		if err != nil {
			t.Fatalf("Failed to write response: %v", err)
		}
	}))
}

// TestListAccounts tests the ListAccounts method
func TestListAccounts(t *testing.T) {
	// Create a mock server that returns a JSON response
	mockResp := `{
		"data": [
			{
				"id": "1",
				"name": "Test Account",
				"type": "asset",
				"current_balance": "1000.00",
				"currency_code": "USD"
			}
		]
	}`
	
	server := mockServer(t, http.StatusOK, mockResp)
	defer server.Close()
	
	// Create a client pointing to the mock server
	client := NewClient(server.URL, "test-token")
	
	// Call the method being tested
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	accounts, err := client.ListAccounts(ctx, 1, 10)
	
	// Assert the results
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}
	
	if len(accounts) != 1 {
		t.Fatalf("Expected 1 account, got %d", len(accounts))
	}
	
	account := accounts[0]
	if account.ID != "1" {
		t.Errorf("Expected ID '1', got '%s'", account.ID)
	}
	
	if account.Name != "Test Account" {
		t.Errorf("Expected Name 'Test Account', got '%s'", account.Name)
	}
	
	// Add more assertions as needed
}

// TestErrorHandling tests how the client handles API errors
func TestErrorHandling(t *testing.T) {
	// Create a mock server that returns an error response
	mockResp := `{
		"message": "Resource not found",
		"errors": {
			"id": ["Invalid ID provided"]
		}
	}`
	
	server := mockServer(t, http.StatusNotFound, mockResp)
	defer server.Close()
	
	// Create a client pointing to the mock server
	client := NewClient(server.URL, "test-token")
	
	// Call a method that should fail
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	_, err := client.GetAccount(ctx, "invalid-id")
	
	// Assert that an error was returned
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
	
	// Check that it's the correct type of error
	// This will depend on your error handling implementation
	// For example:
	// apiErr, ok := err.(*APIError)
	// if !ok {
	//     t.Fatalf("Expected APIError, got %T", err)
	// }
	//
	// if apiErr.StatusCode != http.StatusNotFound {
	//     t.Errorf("Expected status code %d, got %d", http.StatusNotFound, apiErr.StatusCode)
	// }
} 