package cmd

import (
	"fmt"
	"log"
	"strings"
	"os"
	"path/filepath"

	ecrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/spf13/cobra"
	"github.com/trigg3rX/triggerx-keeper/pkg/core"
	"github.com/trigg3rX/triggerx-keeper/pkg/core/crypto"
	// "github.com/trigg3rX/triggerx-keeper/pkg/core/logger"
	"github.com/trigg3rX/triggerx-avs-cli/utils"
)

var generateKeystoreCmd = &cobra.Command{
	Use:   "generate-keystore",
	Short: "Generate keystore",
	Long: `
	Generate a new keystore for the operator.
	The input required is a password for the new keystore. 
	Password must contain:
	- At least 8 characters
	- At least one uppercase letter
	- At least one lowercase letter
	- At least one number
	- At least one special character
	The keystore file will be encrypted with the provided password. The same password will be used to decrypt the keystore file when signing transactions.
	Keep the password safe as it is required to decrypt the keystore file when signing transactions.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		generateKeystore(cmd)
	},
}

func generateKeystore(cmd *cobra.Command) {
	fmt.Println(`
	Generate a new keystore for the operator.
	The input required is a password for the new keystore. 
	Password must contain:
		- At least 8 characters
		- At least one uppercase letter
		- At least one lowercase letter
		- At least one number
		- At least one special character
	The keystore file will be encrypted with the provided password. The same password will be used to decrypt the keystore file when signing transactions.
	Keep the password safe as it is required to decrypt the keystore file when signing transactions.
	`)

	err := generateAndSaveKeystore()
	if err != nil {
		log.Fatalf("Failed to generate and save keystore: %v", err)
	}
}

func generateAndSaveKeystore() error {
	password, err := utils.PasswordPrompt("Enter password for the new keystore: ", true)
	if err != nil {
		return fmt.Errorf("failed to get password: %w", err)
	}

	keyType, err := utils.StringPrompt("Enter key type (ecdsa/bls): ", false, true, []string{"ecdsa", "bls"})
	if err != nil {
		return fmt.Errorf("failed to get key type: %w", err)
	}

	keyType = strings.ToLower(strings.TrimSpace(keyType))
	if keyType != "ecdsa" && keyType != "bls" {
		log.Fatalf("Invalid key type. Must be either 'ecdsa' or 'bls'")
	}

	switch keyType {
	case "ecdsa":
		privateKey, err := ecrypto.GenerateKey()
		if err != nil {
			return fmt.Errorf("failed to generate ECDSA key: %w", err)
		}
		privKey := ecrypto.FromECDSA(privateKey)
		addr := ecrypto.PubkeyToAddress(privateKey.PublicKey)
		return saveKeystore(keyType, password, addr.Bytes(), privKey)
	case "bls":
		blsScheme := crypto.NewBLSScheme(crypto.BN254)
		privKey, err := blsScheme.GenerateRandomKey()
		if err != nil {
			return fmt.Errorf("failed to generate BLS key: %w", err)
		}
		pubKey, err := blsScheme.GetPublicKey(privKey, false, true)
		if err != nil {
			return fmt.Errorf("failed to get BLS public key: %w", err)
		}
		return saveKeystore(keyType, password, pubKey, privKey)
	default:
		return fmt.Errorf("invalid key type: %s", keyType)
	}

}

func saveKeystore(keyType string, password string, pubKey, privKey []byte) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}
	ksPath := filepath.Join(homeDir, ".triggerx/keystore", fmt.Sprintf("%s_%x.key", keyType, pubKey[:6]))
	var cryptoCurve crypto.CryptoCurve
	if keyType == "ecdsa" {
		cryptoCurve = crypto.CryptoCurve("ECDSA")
	} else if keyType == "bls" {
		cryptoCurve = crypto.CryptoCurve("BN254")
	}
	if err := utils.DisplayWarningMessage(keyType, core.Bytes2Hex(privKey), ksPath); err != nil {
		return fmt.Errorf("failed to display warning message: %w", err)
	}
	return crypto.SaveKey(cryptoCurve, privKey, password, ksPath)
}