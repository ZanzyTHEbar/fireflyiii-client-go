package commands

import (
	"fmt"
	"log"

	firefly "github.com/ZanzyTHEbar/fireflyiii-client-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// accountsCmd represents the accounts command
var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Manage Firefly III accounts",
	Long: `List and manage accounts in your Firefly III instance.
	
Examples:
  firefly-client accounts list
  firefly-client accounts list --type=asset
  firefly-client accounts show 123`,
}

var accountsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts",
	Long:  `List all accounts from your Firefly III instance`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Implement account listing
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

		// TODO: Implement actual account listing using the client
		// Example: accounts, err := client.GetAccounts(context.Background())
		fmt.Println("TODO: Implement account listing with the Firefly client")

		// For now, just show that the CLI structure works
		fmt.Println("Accounts command is working! Implementation coming soon...")
	},
}

var accountsShowCmd = &cobra.Command{
	Use:   "show [account-id]",
	Short: "Show details of a specific account",
	Long:  `Show detailed information about a specific account by ID`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		accountID := args[0]
		// TODO: Implement account details
		fmt.Printf("TODO: Show details for account ID: %s\n", accountID)
	},
}

func init() {
	rootCmd.AddCommand(accountsCmd)
	accountsCmd.AddCommand(accountsListCmd)
	accountsCmd.AddCommand(accountsShowCmd)

	// TODO: Add flags for account type filtering, pagination, etc.
	// accountsListCmd.Flags().String("type", "", "Filter by account type (asset, expense, revenue, etc.)")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
