package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var fireflyURL string
var token string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "firefly-client",
	Short: "A CLI client for Firefly III",
	Long: `A command-line interface to interact with your Firefly III instance.
Complete documentation is available at https://github.com/ZanzyTHEbar/firefly-client-go`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.firefly-client/config.yaml or ./config.yaml)")
	rootCmd.PersistentFlags().StringVarP(&fireflyURL, "url", "u", "", "Firefly III instance URL (e.g., http://localhost:8080)")
	rootCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "Firefly III API token")

	// TODO: Add more persistent flags as needed, e.g., for output format (json, yaml, text)

	vipErr := viper.BindPFlag("firefly_url", rootCmd.PersistentFlags().Lookup("url"))
	cobra.CheckErr(vipErr)
	vipErr = viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	cobra.CheckErr(vipErr)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".firefly-client" (without extension).
		viper.AddConfigPath(home + "/.firefly-client")
		viper.AddConfigPath(".") // Also look in the current directory
		viper.SetConfigName("config")
		viper.SetConfigType("yaml") // Can be yaml, json, toml, etc.
	}

	viper.SetEnvPrefix("FIREFLY") // Will look for FIREFLY_URL, FIREFLY_TOKEN
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if it's just not found
			// fmt.Fprintln(os.Stderr, "Config file not found, relying on flags and environment variables.")
		} else {
			// Config file was found but another error was produced
			fmt.Fprintln(os.Stderr, "Error reading config file:", err)
		}
	}

	// TODO: Add validation for required config values (URL, Token)
	// For example:
	// if viper.GetString("firefly_url") == "" {
	// 	fmt.Fprintln(os.Stderr, "Error: firefly_url is not set. Please set it via config file, environment variable FIREFLY_URL, or --url flag.")
	// 	os.Exit(1)
	// }
	// if viper.GetString("token") == "" {
	//  fmt.Fprintln(os.Stderr, "Error: token is not set. Please set it via config file, environment variable FIREFLY_TOKEN, or --token flag.")
	// 	os.Exit(1)
	// }
}
