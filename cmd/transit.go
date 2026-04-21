package cmd

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpipe/internal/vault"
)

var (
	transitKey     string
	transitDecrypt bool
)

func init() {
	transitCmd := &cobra.Command{
		Use:   "transit <plaintext|ciphertext>",
		Short: "Encrypt or decrypt data using Vault Transit",
		Args:  cobra.ExactArgs(1),
		RunE:  runTransit,
	}
	transitCmd.Flags().StringVar(&transitKey, "key", "", "Transit key name (required)")
	transitCmd.Flags().BoolVar(&transitDecrypt, "decrypt", false, "Decrypt instead of encrypt")
	_ = transitCmd.MarkFlagRequired("key")
	rootCmd.AddCommand(transitCmd)
}

func runTransit(cmd *cobra.Command, args []string) error {
	c, err := vault.NewClient("", "")
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	if transitDecrypt {
		res, err := vault.DecryptTransit(c, transitKey, args[0])
		if err != nil {
			return fmt.Errorf("decrypt: %w", err)
		}
		decoded, err := base64.StdEncoding.DecodeString(res.Plaintext)
		if err != nil {
			// Return raw base64 if caller wants it undecoded
			fmt.Fprintln(os.Stdout, res.Plaintext)
			return nil
		}
		fmt.Fprintln(os.Stdout, string(decoded))
		return nil
	}

	// Encrypt: base64-encode the input first
	b64 := base64.StdEncoding.EncodeToString([]byte(args[0]))
	res, err := vault.EncryptTransit(c, transitKey, b64)
	if err != nil {
		return fmt.Errorf("encrypt: %w", err)
	}
	fmt.Fprintln(os.Stdout, res.Ciphertext)
	return nil
}
