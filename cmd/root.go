package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trigg3rX/triggerx-avs-cli/utils"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "triggerx",
	Short: "TriggerX AVS CLI",
	Long:  utils.TriggerXAsciiArt + "\nTriggerX AVS CLI - Tool for managing operators on AVS.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	registerCmd.GroupID = "1_core"
	deregisterCmd.GroupID = "1_core"
	statusCmd.GroupID = "1_core"

	generateKeystoreCmd.GroupID = "2_keystore"

	versionCmd.GroupID = "5_utils"

	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(deregisterCmd)
	rootCmd.AddCommand(statusCmd)

	rootCmd.AddCommand(generateKeystoreCmd)

	rootCmd.AddCommand(versionCmd)

	rootCmd.AddGroup(&cobra.Group{
		ID:    "1_core",
		Title: "Core Commands:",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "2_keystore",
		Title: "Keystore Commands:",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "5_utils",
		Title: "Utility Commands:",
	})

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".triggerx")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
