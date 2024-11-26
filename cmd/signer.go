package cmd

import (
	"github.com/spf13/cobra"
)

var updateSignerCmd = &cobra.Command{
	Use:   "update-signer",
	Short: "Update signer address at provided index",
	Run: func(cmd *cobra.Command, args []string) {
		updateSigner(cmd)
	},
}

var removeSignerCmd = &cobra.Command{
	Use:   "remove-signer",
	Short: "Remove signer address at provided index",
	Run: func(cmd *cobra.Command, args []string) {
		removeSigner(cmd)
	},
}

var addSignerCmd = &cobra.Command{
	Use:   "add-signer",
	Short: "Add signer address",
	Run: func(cmd *cobra.Command, args []string) {
		addSigner(cmd)
	},
}

var listSignersCmd = &cobra.Command{
	Use:   "list-signers",
	Short: "List signer addresses",
	Run: func(cmd *cobra.Command, args []string) {
		listSigners(cmd)
	},
}

func updateSigner(cmd *cobra.Command) {
	// Implementation goes here
}

func removeSigner(cmd *cobra.Command) {
	// Implementation goes here
}

func addSigner(cmd *cobra.Command) {
	// Implementation goes here
}

func listSigners(cmd *cobra.Command) {
	// Implementation goes here
}
