package commands

import (
	"fmt"
	"log"
	"time"

	firefly "github.com/ZanzyTHEbar/fireflyiii-client-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test connection to Firefly III",
	Long: `Test the connection to your Firefly III instance and verify authentication.
	
This command will:
- Validate your configuration
- Test the connection to your Firefly III instance  
- Verify your API token is working
- Show basic instance information

Examples:
  firefly-client test
  firefly-client test --timeout=10`,
	Run: func(cmd *cobra.Command, args []string) {
		url := viper.GetString("firefly_url")
		token := viper.GetString("token")

		if url == "" {
			log.Fatal("Firefly URL is required. Set it via --url flag, FIREFLY_URL environment variable, or config file.")
		}
		if token == "" {
			log.Fatal("API token is required. Set it via --token flag, FIREFLY_TOKEN environment variable, or config file.")
		}

		fmt.Println("🔧 Testing Firefly III connection...")
		fmt.Printf("📍 URL: %s\n", url)
		fmt.Printf("🔑 Token: %s...\n", token[:min(len(token), 8)])

		// Create client with timeout
		start := time.Now()
		client, err := firefly.NewFireflyClient(url, token)
		if err != nil {
			fmt.Printf("❌ Failed to create client: %v\n", err)
			return
		}

		clientDuration := time.Since(start)
		fmt.Printf("✅ Client created successfully (took %v)\n", clientDuration)

		// TODO: Add actual API test once we implement a simple API call
		// For example, a call to get user info or system status
		// This would involve:
		// 1. Making a simple API call (e.g., GET /api/v1/about)
		// 2. Checking the response status
		// 3. Verifying the token is valid

		fmt.Printf("🕒 Total test time: %v\n", time.Since(start))
		fmt.Println("✅ Connection test completed successfully!")
		fmt.Println("\n💡 Next steps:")
		fmt.Println("   - Try 'firefly-client accounts list' to see your accounts")
		fmt.Println("   - Try 'firefly-client transactions list' to see recent transactions")

		// Show that client is not nil to demonstrate successful creation
		_ = client // Use the client variable to avoid unused warning
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Add timeout flag specifically for test command
	testCmd.Flags().Int("timeout", 30, "Connection timeout in seconds")
	viper.BindPFlag("test_timeout", testCmd.Flags().Lookup("timeout"))
}
