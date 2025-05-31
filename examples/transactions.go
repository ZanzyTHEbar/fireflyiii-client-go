package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	firefly "github.com/ZanzyTHEbar/fireflyiii-client-go"
)

func main() {
	// Get Firefly III API URL and personal access token from environment variables
	baseURL := os.Getenv("FIREFLY_BASE_URL")
	if baseURL == "" {
		log.Fatal("FIREFLY_BASE_URL environment variable is required")
	}

	token := os.Getenv("FIREFLY_ACCESS_TOKEN")
	if token == "" {
		log.Fatal("FIREFLY_ACCESS_TOKEN environment variable is required")
	}

	// Initialize client
	client, err := firefly.NewClient(baseURL, token)
	if err != nil {
		log.Fatalf("Error initializing client: %v", err)
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get a list of accounts to use for transactions
	accounts, err := client.ListAccounts(ctx, 1, 5)
	if err != nil {
		log.Fatalf("Error listing accounts: %v", err)
	}

	if len(accounts) == 0 {
		log.Fatal("No accounts found. Please create an account first.")
	}

	sourceAccount := accounts[0]
	fmt.Printf("Using account: %s (ID: %s) for transactions\n\n", sourceAccount.Name, sourceAccount.ID)

	// Create a new transaction
	fmt.Println("Creating a new transaction:")
	tx := firefly.TransactionModel{
		Description: "Example transaction from Go client",
		Amount:      42.50,
		Date:        time.Now().Format("2006-01-02"),
		AccountID:   sourceAccount.ID,
		Type:        "withdrawal",
		CategoryName: "Groceries",
		Notes:       "This is a test transaction created by the example code",
	}

	err = client.ImportTransaction(ctx, tx)
	if err != nil {
		log.Fatalf("Error creating transaction: %v", err)
	}
	fmt.Println("Transaction created successfully!")

	// List transactions
	fmt.Println("\nListing recent transactions:")
	transactions, err := client.ListTransactions(ctx, 1, 5)
	if err != nil {
		log.Fatalf("Error listing transactions: %v", err)
	}

	for i, tx := range transactions {
		fmt.Printf("%d. %s - %s %s (%s)\n",
			i+1,
			tx.Date,
			tx.Amount,
			tx.CurrencyCode,
			tx.Description,
		)
	}

	// Search for transactions
	fmt.Println("\nSearching for transactions containing 'example':")
	searchResults, err := client.SearchTransactions(ctx, "example")
	if err != nil {
		log.Fatalf("Error searching transactions: %v", err)
	}

	if len(searchResults) == 0 {
		fmt.Println("No transactions found with that search term.")
		return
	}

	// Get the first transaction from search results
	ourTransaction := searchResults[0]
	fmt.Printf("Found transaction: %s - %s (%s)\n",
		ourTransaction.Date,
		ourTransaction.Amount,
		ourTransaction.Description,
	)

	// Get transaction details
	fmt.Println("\nFetching transaction details:")
	txDetails, err := client.GetTransaction(ctx, ourTransaction.ID)
	if err != nil {
		log.Fatalf("Error getting transaction details: %v", err)
	}

	fmt.Printf("Transaction details:\n")
	fmt.Printf("  Date: %s\n", txDetails.Date)
	fmt.Printf("  Description: %s\n", txDetails.Description)
	fmt.Printf("  Amount: %s %s\n", txDetails.Amount, txDetails.CurrencyCode)
	fmt.Printf("  Category: %s\n", txDetails.CategoryName)
	fmt.Printf("  Notes: %s\n", txDetails.Notes)

	// Update the transaction
	fmt.Println("\nUpdating transaction:")
	txDetails.Description = "Updated example transaction"
	txDetails.Amount = "50.00"
	txDetails.CategoryName = "Entertainment"

	err = client.UpdateTransaction(ctx, txDetails.ID, *txDetails)
	if err != nil {
		log.Fatalf("Error updating transaction: %v", err)
	}
	fmt.Println("Transaction updated successfully!")

	// Get updated transaction
	fmt.Println("\nFetching updated transaction details:")
	updatedTx, err := client.GetTransaction(ctx, txDetails.ID)
	if err != nil {
		log.Fatalf("Error getting updated transaction: %v", err)
	}

	fmt.Printf("Updated transaction details:\n")
	fmt.Printf("  Date: %s\n", updatedTx.Date)
	fmt.Printf("  Description: %s\n", updatedTx.Description)
	fmt.Printf("  Amount: %s %s\n", updatedTx.Amount, updatedTx.CurrencyCode)
	fmt.Printf("  Category: %s\n", updatedTx.CategoryName)

	// Delete the transaction (uncomment if you want to actually delete)
	/*
		fmt.Println("\nDeleting the transaction:")
		err = client.DeleteTransaction(ctx, updatedTx.ID)
		if err != nil {
			log.Fatalf("Error deleting transaction: %v", err)
		}
		fmt.Println("Transaction deleted successfully!")
	*/
} 