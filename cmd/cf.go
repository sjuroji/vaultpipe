package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourusername/vaultpipe/internal/vault"
)

var (
	cfRole              string
	cfCertFile          string
	cfKeyFile           string
	cfSigningTime       string
	cfSignature         string
	cfMount             string
)

func init() {
	cfCmd := &cobra.Command{
		Use:   "cf",
		Short: "Authenticate using Cloud Foundry auth method",
		RunE:  runCFLogin,
	}
	cfCmd.Flags().StringVar(&cfRole, "role", "", "CF role name (required)")
	cfCmd.Flags().StringVar(&cfCertFile, "cert-file", "", "Path to CF instance cert file (required)")
	cfCmd.Flags().StringVar(&cfKeyFile, "key-file", "", "Path to CF instance key file (required)")
	cfCmd.Flags().StringVar(&cfSigningTime, "signing-time", "", "Signing time in RFC3339 format (required)")
	cfCmd.Flags().StringVar(&cfSignature, "signature", "", "Request signature (required)")
	cfCmd.Flags().StringVar(&cfMount, "mount", "cf", "Auth mount path")
	_ = cfCmd.MarkFlagRequired("role")
	_ = cfCmd.MarkFlagRequired("cert-file")
	_ = cfCmd.MarkFlagRequired("key-file")
	_ = cfCmd.MarkFlagRequired("signing-time")
	_ = cfCmd.MarkFlagRequired("signature")
	rootCmd.AddCommand(cfCmd)
}

func runCFLogin(cmd *cobra.Command, args []string) error {
	certBytes, err := os.ReadFile(cfCertFile)
	if err != nil {
		return fmt.Errorf("read cert file: %w", err)
	}
	keyBytes, err := os.ReadFile(cfKeyFile)
	if err != nil {
		return fmt.Errorf("read key file: %w", err)
	}

	client, err := vault.NewClient("", "")
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	token, err := vault.CFLogin(client, cfRole, string(certBytes), string(keyBytes), cfSigningTime, cfSignature, cfMount)
	if err != nil {
		return fmt.Errorf("cf login: %w", err)
	}

	export, _ := cmd.Flags().GetBool("export")
	if export {
		fmt.Printf("export VAULT_TOKEN=%s\n", token)
	} else {
		fmt.Println(token)
	}
	return nil
}
