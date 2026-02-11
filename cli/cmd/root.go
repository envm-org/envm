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
	_ = godotenv.Load()

	viper.SetEnvPrefix("ENVM_CLI")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

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
