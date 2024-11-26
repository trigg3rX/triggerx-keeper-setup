package cmd

import (
	"github.com/spf13/cobra"
)

var updateBLSCmd = &cobra.Command{
	Use:   "update-bls",
	Short: "Update BLS public key at provided index",
	Run: func(cmd *cobra.Command, args []string) {
		updateBLS(cmd)
	},
}

var removeBLSCmd = &cobra.Command{
	Use:   "remove-bls",
	Short: "Remove BLS public key at provided index",
	Run: func(cmd *cobra.Command, args []string) {
		removeBLS(cmd)
	},
}

var addBLSCmd = &cobra.Command{
	Use:   "add-bls",
	Short: "Add BLS public key",
	Run: func(cmd *cobra.Command, args []string) {
		addBLS(cmd)
	},
}

var listBLSCmd = &cobra.Command{
	Use:   "list-bls",
	Short: "List BLS public keys",
	Run: func(cmd *cobra.Command, args []string) {
		listBLS(cmd)
	},
}

func updateBLS(cmd *cobra.Command) {
	// Implementation goes here
}

func removeBLS(cmd *cobra.Command) {
	// Implementation goes here
}

func addBLS(cmd *cobra.Command) {
	// Implementation goes here
}

func listBLS(cmd *cobra.Command) {
	// Implementation goes here
}
