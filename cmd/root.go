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

	// registerCmd.GroupID = "1_core"
	// deregisterCmd.GroupID = "1_core"
	statusCmd.GroupID = "1_core"

	generateKeystoreCmd.GroupID = "2_keystore"
	// exportKeystoreCmd.GroupID = "2_keystore"
	// importKeystoreCmd.GroupID = "2_keystore"

	// updateBLSCmd.GroupID = "3_bls"
	// removeBLSCmd.GroupID = "3_bls"
	// addBLSCmd.GroupID = "3_bls"
	// listBLSCmd.GroupID = "3_bls"

	// updateSignerCmd.GroupID = "4_signer"
	// removeSignerCmd.GroupID = "4_signer"
	// addSignerCmd.GroupID = "4_signer"
	// listSignersCmd.GroupID = "4_signer"

	versionCmd.GroupID = "5_utils"

	// rootCmd.AddCommand(registerCmd)
	// rootCmd.AddCommand(deregisterCmd)
	rootCmd.AddCommand(statusCmd)

	rootCmd.AddCommand(generateKeystoreCmd)
	// rootCmd.AddCommand(exportKeystoreCmd)
	// rootCmd.AddCommand(importKeystoreCmd)

	// rootCmd.AddCommand(updateBLSCmd)
	// rootCmd.AddCommand(removeBLSCmd)
	// rootCmd.AddCommand(addBLSCmd)
	// rootCmd.AddCommand(listBLSCmd)

	// rootCmd.AddCommand(updateSignerCmd)
	// rootCmd.AddCommand(removeSignerCmd)
	// rootCmd.AddCommand(addSignerCmd)
	// rootCmd.AddCommand(listSignersCmd)

	rootCmd.AddCommand(versionCmd)

	rootCmd.AddGroup(&cobra.Group{
		ID:    "1_core",
		Title: "Core Commands:",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "2_keystore",
		Title: "Keystore Commands:",
	})
	// rootCmd.AddGroup(&cobra.Group{
	// 	ID:    "3_bls",
	// 	Title: "BLS Commands:",
	// })
	// rootCmd.AddGroup(&cobra.Group{
	// 	ID:    "4_signer",
	// 	Title: "Signer Commands:",
	// })
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
		viper.SetConfigName(".triggerx-avs-cli")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
