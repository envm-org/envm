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

// DefaultAPIURL and Version are overridden at build time via ldflags:
//
//	-X 'github.com/envm-org/cli/cmd.DefaultAPIURL=https://envm.onrender.com'
//	-X 'github.com/envm-org/cli/cmd.Version=0.0.1'
var (
	cfgFile string

	DefaultAPIURL = "http://localhost:8080"
	Version       = "dev"

	rootCmd = &cobra.Command{
		Use:   "envm",
		Short: "Command line interface for envm",
		// PersistentPreRunE is called after flags are parsed but before the
		// command's RunE function is called.
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
)

func Execute() {
	// Set version here so the ldflags-injected value is always picked up.
	// Assigning in the struct literal captures the zero-value before ldflags runs.
	rootCmd.Version = Version
	ui.SetupCobraHelp(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default locations: ., $HOME/.envm-cli/)")
	rootCmd.PersistentFlags().String("api-url", DefaultAPIURL, "API URL (default: "+DefaultAPIURL+")")
	viper.BindPFlag("api-url", rootCmd.PersistentFlags().Lookup("api-url"))

	// Hide global flags from the help display
	rootCmd.PersistentFlags().MarkHidden("config")
	rootCmd.PersistentFlags().MarkHidden("api-url")
}

func initializeConfig(cmd *cobra.Command) error {
	_ = godotenv.Load()

	viper.SetEnvPrefix("ENVM_CLI")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.BindEnv("api-url", "API_URL")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		// Only panic if we can't get the home directory.
		cobra.CheckErr(err)

		viper.AddConfigPath(".")
		viper.AddConfigPath(home + "/.envm-cli")
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
	}

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}

	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return err
	}

	return nil
}
