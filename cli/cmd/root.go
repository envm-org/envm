package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/envm-org/cli/internal/ui"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "cli",
		Short: "Command line interface for envm",
		// PersistentPreRunE is called after flags are parsed but before the
		// command's RunE function is called.
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			ui.PrintLogo()
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default locations: ., $HOME/.envm-cli/)")
	rootCmd.PersistentFlags().String("api-url", "http://localhost:8080", "API URL (default: http://localhost:8080)")
	viper.BindPFlag("api-url", rootCmd.PersistentFlags().Lookup("api-url"))
}

func initializeConfig(cmd *cobra.Command) error {
	// Load .env file if it exists
	_ = godotenv.Load()

	// 1. Set up Viper to use environment variables.
	viper.SetEnvPrefix("ENVM_CLI")
	// Allow for nested keys in environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// 2. Handle the configuration file.
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search for a config file in default locations.
		home, err := os.UserHomeDir()
		// Only panic if we can't get the home directory.
		cobra.CheckErr(err)

		// Search for a config file with the name "config" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home + "/.envm-cli")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	// 3. Read the configuration file.
	// If a config file is found, read it in. We use a robust error check
	// to ignore "file not found" errors, but panic on any other error.
	if err := viper.ReadInConfig(); err != nil {
		// It's okay if the config file doesn't exist.
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}

	// 4. Bind Cobra flags to Viper.
	// This is the magic that makes the flag values available through Viper.
	// It binds the full flag set of the command passed in.
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return err
	}

	return nil
}
