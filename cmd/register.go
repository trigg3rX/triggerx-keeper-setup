package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	// "github.com/spf13/viper"
	"github.com/spf13/cobra"

	// "github.com/consensys/gnark-crypto/ecc/bn254"
	// ecommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

// 	"github.com/Layr-Labs/eigensdk-go/chainio/clients/elcontracts"
// 	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
// 	"github.com/Layr-Labs/eigensdk-go/chainio/clients/wallet"
// 	"github.com/Layr-Labs/eigensdk-go/chainio/txmgr"
// 	"github.com/Layr-Labs/eigensdk-go/signerv2"
	eigensdkbls "github.com/Layr-Labs/eigensdk-go/crypto/bls"
	eigensdkecdsa "github.com/Layr-Labs/eigensdk-go/crypto/ecdsa"
	
	// "github.com/trigg3rX/triggerx-keeper/config"

	"github.com/trigg3rX/triggerx-avs-cli/utils"

	// regcoord "github.com/trigg3rX/go-backend/pkg/avsinterface/bindings/RegistryCoordinator"
	// stakeregistry "github.com/trigg3rX/go-backend/pkg/avsinterface/bindings/StakeRegistry"
	// triggerxsm "github.com/trigg3rX/go-backend/pkg/avsinterface/bindings/TriggerXServiceManager"
)

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new operator",
	Long: `
	Register a new operator to the TriggerX AVS. It is assumed that the operator has already been registered to the EigenLayer protocol.
	If not, please register to EigenLayer first. Follow the instructions here: https://github.com/Layr-Labs/eigenlayer-cli/blob/master/README.md

	Keystore files (ECDSA and BLS) are required to be present in the keystore directory, at the path: ~/.triggerx/keystore/<keystore-file>.key 
	The passphrase for the keystore files is required to be entered interactively.

	If no keystore files are present, please generate the requiredkeystore files using the 'generate-keystore' command.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		registerOperator(cmd)
	},
}

var deregisterCmd = &cobra.Command{
	Use:   "deregister",
	Short: "Deregister an operator",
	Run: func(cmd *cobra.Command, args []string) {
		deregisterOperator(cmd)
	},
}

func registerOperator(cmd *cobra.Command) {
	fmt.Println(`
Register a new operator to the TriggerX AVS. It is assumed that the operator has already been registered to the EigenLayer protocol.
	If not, please register to EigenLayer first. Follow the instructions here: https://github.com/Layr-Labs/eigenlayer-cli/blob/master/README.md

	Keystore files (ECDSA and BLS) are required to be present in the keystore directory, at the path: ~/.triggerx/keystore/<keystore-file>.key 
	The passphrase for the keystore files is required to be entered interactively.

	If no keystore files are present, please generate the requiredkeystore files using the 'generate-keystore' command.
	`)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		cmd.PrintErrf("error getting home directory: %v\n", err)
		return
	}

	keystoreDir := filepath.Join(homeDir, ".triggerx", "keystore")

	files, err1 := os.ReadDir(keystoreDir)
	keystoreFiles := []string{}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".key") {
			keystoreFiles = append(keystoreFiles, file.Name())
		}
	}

	if len(keystoreFiles) == 0 || err1 != nil {
		cmd.PrintErrf("Keystore files not found in %s\n", keystoreDir)
		return
	}

	fmt.Println("Found keystore files:")
	for i, file := range keystoreFiles {
		fmt.Printf("%d. %s\n", i+1, file)
	}
	fmt.Print("Enter the BLS and ECDSA keystore file numbers (separated by space): ")
	var blsIndex, ecdsaIndex int
	fmt.Scanf("%d %d", &blsIndex, &ecdsaIndex)
	if blsIndex < 1 || blsIndex > len(keystoreFiles) || ecdsaIndex < 1 || ecdsaIndex > len(keystoreFiles) {
		cmd.PrintErrf("Invalid selection. Please choose numbers between 1 and %d\n", len(keystoreFiles))
		return
	}
	blsFile := keystoreFiles[blsIndex-1]
	ecdsaFile := keystoreFiles[ecdsaIndex-1]

	err = registerOperatorToTriggerX(blsFile, ecdsaFile)
	if err != nil {
		cmd.PrintErrf("error registering operator: %v\n", err)
		return
	}
}

func registerOperatorToTriggerX(blsFile, ecdsaFile string) error {
	// Get passphrases
	blsPassphrase, err := utils.PasswordPrompt("Enter passphrase for the bls keystore file: ", false)
	if err != nil {
		return fmt.Errorf("failed to get bls passphrase: %w", err)
	}

	ecdsaPassphrase, err := utils.PasswordPrompt("Enter passphrase for the ecdsa keystore file: ", false)
	if err != nil {
		return fmt.Errorf("failed to get ecdsa passphrase: %w", err)
	}

	// Get home directory and construct full paths
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}
	keystoreDir := filepath.Join(homeDir, ".triggerx", "keystore")
	blsKeyPath := filepath.Join(keystoreDir, blsFile)
	ecdsaKeyPath := filepath.Join(keystoreDir, ecdsaFile)

	// Read and decrypt BLS key
	blsKeyPair, err := eigensdkbls.ReadPrivateKeyFromFile(blsKeyPath, blsPassphrase)
	if err != nil {
		return fmt.Errorf("failed to read BLS key: %w", err)
	}

	// Read and decrypt ECDSA key
	ecdsaPrivKey, err := eigensdkecdsa.ReadKey(ecdsaKeyPath, ecdsaPassphrase)
	if err != nil {
		return fmt.Errorf("failed to read ECDSA key: %w", err)
	}

	// Get the operator address from ECDSA key
	operatorAddr := crypto.PubkeyToAddress(ecdsaPrivKey.PublicKey)
	
	// Get BLS public key in G1 and G2 forms
	blsPubKeyG1 := blsKeyPair.GetPubKeyG1()
	blsPubKeyG2 := blsKeyPair.GetPubKeyG2()

	fmt.Printf("Successfully loaded keys for operator address: %s\n", operatorAddr.Hex())
	fmt.Printf("BLS Public Key G1: %x\n", blsPubKeyG1.Serialize())
	fmt.Printf("BLS Public Key G2: %x\n", blsPubKeyG2.Serialize())

	fmt.Printf("ECDSA Public Key: %x\n", ecdsaPrivKey.PublicKey.X.Bytes())
	fmt.Printf("ECDSA Public Key Y: %x\n", ecdsaPrivKey.PublicKey.Y.Bytes())

	fmt.Printf("BLS Private Key: %x\n", blsKeyPair.PrivKey.Bytes())
	fmt.Printf("ECDSA Private Key: %x\n", crypto.FromECDSA(ecdsaPrivKey))
	// TODO: Use these keys to register the operator with the AVS
	// This would typically involve creating and sending transactions
	// to the appropriate smart contracts

	return nil
}

func deregisterOperator(cmd *cobra.Command) {
	// Implementation goes here
}