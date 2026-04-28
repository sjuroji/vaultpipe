package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/vault"
)

var (
	cfMount       string
	cfRole        string
	cfSigningTime string
	cfCert        string
	cfSignature   string
	cfExport      bool
)

func init() {
	cfCmd := &cobra.Command{
		Use:   "cf",
		Short: "Authenticate using the Cloud Foundry auth method",
		RunE:  runCFLogin,
	}

	cfCmd.Flags().StringVar(&cfMount, "mount", "", "CF auth mount path (default: cf)")
	cfCmd.Flags().StringVar(&cfRole, "role", "", "CF role name (required)")
	cfCmd.Flags().StringVar(&cfSigningTime, "signing-time", "", "Signing time in RFC3339 format (required)")
	cfCmd.Flags().StringVar(&cfCert, "cert", "", "CF instance certificate (required)")
	cfCmd.Flags().StringVar(&cfSignature, "signature", "", "Request signature (required)")
	cfCmd.Flags().BoolVar(&cfExport, "export", false, "Print token as export statement")

	_ = cfCmd.MarkFlagRequired("role")
	_ = cfCmd.MarkFlagRequired("signing-time")
	_ = cfCmd.MarkFlagRequired("cert")
	_ = cfCmd.MarkFlagRequired("signature")

	rootCmd.AddCommand(cfCmd)
}

func runCFLogin(cmd *cobra.Command, args []string) error {
	client, err := vault.NewClient("", "")
	if err != nil {
		return fmt.Errorf("cf login: %w", err)
	}

	token, err := vault.CFLogin(client, cfMount, cfRole, cfSigningTime, cfCert, cfSignature)
	if err != nil {
		return fmt.Errorf("cf login: %w", err)
	}

	if cfExport {
		fmt.Fprintf(os.Stdout, "export VAULT_TOKEN=%s\n", token)
	} else {
		fmt.Fprintln(os.Stdout, token)
	}
	return nil
}
