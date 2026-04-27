package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/vault"
)

var (
	cfRole     string
	cfMount    string
	cfExport   bool
)

func init() {
	cfCmd := &cobra.Command{
		Use:   "cf",
		Short: "Authenticate with Vault using Cloud Foundry",
		RunE:  runCFLogin,
	}

	cfCmd.Flags().StringVar(&cfRole, "role", "", "CF role name (required)")
	cfCmd.Flags().StringVar(&cfMount, "mount", "cf", "Auth mount path")
	cfCmd.Flags().BoolVar(&cfExport, "export", false, "Print token as export statement")
	_ = cfCmd.MarkFlagRequired("role")

	rootCmd.AddCommand(cfCmd)
}

func runCFLogin(cmd *cobra.Command, args []string) error {
	signingTime := os.Getenv("CF_SIGNING_TIME")
	instanceCert := os.Getenv("CF_INSTANCE_CERT")
	signature := os.Getenv("CF_SIGNATURE")

	client, err := vault.NewClient("", "")
	if err != nil {
		return fmt.Errorf("cf login: create client: %w", err)
	}

	token, err := vault.CFLogin(client, vault.CFLoginRequest{
		RoleID:         cfRole,
		SigningTime:    signingTime,
		CFInstanceCert: instanceCert,
		Signature:      signature,
		Mount:          cfMount,
	})
	if err != nil {
		return fmt.Errorf("cf login: %w", err)
	}

	if cfExport {
		fmt.Fprintf(cmd.OutOrStdout(), "export VAULT_TOKEN=%s\n", token)
	} else {
		fmt.Fprintln(cmd.OutOrStdout(), token)
	}
	return nil
}
