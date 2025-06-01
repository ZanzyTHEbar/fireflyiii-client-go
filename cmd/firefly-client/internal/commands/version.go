package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version is set at build time
var Version = "dev"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of firefly-client",
	Long:  `All software has versions. This is firefly-client's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("firefly-client version %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
