package commands

import (
	"fmt"
	"log"

	firefly "github.com/ZanzyTHEbar/fireflyiii-client-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// transactionsCmd represents the transactions command
var transactionsCmd = &cobra.Command{
	Use:   "transactions",
	Short: "Manage Firefly III transactions",
	Long: `List and manage transactions in your Firefly III instance.
	
Examples:
  firefly-client transactions list
  firefly-client transactions list --limit=50
  firefly-client transactions show 123`,
}

var transactionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List recent transactions",
	Long:  `List recent transactions from your Firefly III instance`,
	Run: func(cmd *cobra.Command, args []string) {
		url := viper.GetString("firefly_url")
		token := viper.GetString("token")

		if url == "" {
			log.Fatal("Firefly URL is required. Set it via --url flag, FIREFLY_URL environment variable, or config file.")
		}
		if token == "" {
			log.Fatal("API token is required. Set it via --token flag, FIREFLY_TOKEN environment variable, or config file.")
		}

		// Create Firefly client
		client, err := firefly.NewFireflyClient(url, token)
		if err != nil {
			log.Fatalf("Failed to create Firefly client: %v", err)
		}

		fmt.Printf("Connecting to Firefly III at: %s\n", url)
		fmt.Printf("Using token: %s...\n", token[:min(len(token), 8)])
		fmt.Printf("Client created successfully: %v\n", client != nil)

		// TODO: Implement actual transaction listing using the client
		// Example: transactions, err := client.GetTransactions(context.Background())
		fmt.Println("TODO: Implement transaction listing with the Firefly client")

		// For now, just show that the CLI structure works
		fmt.Println("Transactions command is working! Implementation coming soon...")
	},
}

var transactionsShowCmd = &cobra.Command{
	Use:   "show [transaction-id]",
	Short: "Show details of a specific transaction",
	Long:  `Show detailed information about a specific transaction by ID`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		transactionID := args[0]
		// TODO: Implement transaction details
		fmt.Printf("TODO: Show details for transaction ID: %s\n", transactionID)
	},
}

func init() {
	rootCmd.AddCommand(transactionsCmd)
	transactionsCmd.AddCommand(transactionsListCmd)
	transactionsCmd.AddCommand(transactionsShowCmd)

	// TODO: Add flags for filtering, pagination, etc.
	// transactionsListCmd.Flags().Int("limit", 20, "Number of transactions to return")
	// transactionsListCmd.Flags().String("type", "", "Filter by transaction type")
}
