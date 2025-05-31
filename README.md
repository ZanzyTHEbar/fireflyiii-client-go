# Firefly III Go Client

A comprehensive Go client library for the Firefly III API.

## Overview

This library provides a Go client for interacting with [Firefly III](https://www.firefly-iii.org/), a self-hosted financial management application. The client is based on the Firefly III OpenAPI 3.0 specification (currently targeting `firefly-iii-6.2.8`).

## Features

- ✅ Complete API coverage for Firefly III endpoints
- ✅ Idiomatic Go interface for all operations
- ✅ Support for Personal Access Token authentication
- ✅ Comprehensive error handling
- ✅ Pagination support for list operations
- ✅ Context support for request cancellation

## Installation

```bash
go get github.com/ZanzyTHEbar/firefly-client-go
```

## Quick Start

### Initialize the client

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    firefly "github.com/ZanzyTHEbar/firefly-client-go"
)

func main() {
    // Initialize the client with your Firefly III API URL and personal access token
    baseURL := "https://your-firefly-iii-instance/api"
    token := "your-personal-access-token"
    
    client := firefly.NewClient(baseURL, token)
    
    // Use a context for request cancellation and timeouts
    ctx := context.Background()
    
    // Now you can use the client to interact with your Firefly III instance
    // ...
}
```

### Fetch accounts

```go
// List all accounts
accounts, err := client.ListAccounts(ctx, 1, 25) // Page 1, 25 items per page
if err != nil {
    log.Fatalf("Error listing accounts: %v", err)
}

for _, account := range accounts {
    fmt.Printf("Account: %s (ID: %s)\n", account.Name, account.ID)
}

// Get a specific account by ID
account, err := client.GetAccount(ctx, "account-id")
if err != nil {
    log.Fatalf("Error fetching account: %v", err)
}
fmt.Printf("Account details: %+v\n", account)
```

### Create a transaction

```go
// Create a new transaction
tx := firefly.TransactionModel{
    Description: "Grocery shopping",
    Amount:      "42.50",
    Date:        time.Now().Format("2006-01-02"),
    AccountID:   "your-account-id",
    Type:        "withdrawal",
    // Set other fields as needed
}

if err := client.ImportTransaction(ctx, tx); err != nil {
    log.Fatalf("Error creating transaction: %v", err)
}
fmt.Println("Transaction created successfully")
```

## API Coverage

This client provides access to all Firefly III API endpoints:

- Accounts management
- Transactions
- Categories
- Budgets and budget limits
- Tags
- Piggy banks
- Bills
- Rules
- Reports and charts
- Attachments
- Data import/export

## Error Handling

The client provides detailed error information when API requests fail:

```go
if err := client.ImportTransaction(ctx, tx); err != nil {
    // Check for specific error types
    if apiErr, ok := err.(*firefly.APIError); ok {
        fmt.Printf("API error (%d): %s\n", apiErr.StatusCode, apiErr.Message)
    } else if validErr, ok := err.(*firefly.ValidationError); ok {
        fmt.Println("Validation errors:")
        for field, msgs := range validErr.Errors {
            fmt.Printf("  %s: %v\n", field, msgs)
        }
    } else {
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Documentation

For detailed documentation on all available methods and types, please refer to the [GoDoc](https://pkg.go.dev/github.com/ZanzyTHEbar/firefly-client-go).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgements

- [Firefly III](https://www.firefly-iii.org/) - The awesome free and open-source personal finance manager
- [Firefly III API Documentation](https://api-docs.firefly-iii.org/)
