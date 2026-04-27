package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vaultpipe/vaultpipe/internal/vault"
)

var (
	oidcMount string
	oidcRole  string
	oidcExport bool
)

func init() {
	oidcCmd := &cobra.Command{
		Use:   "oidc",
		Short: "Authenticate with Vault using OIDC",
		RunE:  runOIDCLogin,
	}

	oidcCmd.Flags().StringVar(&oidcMount, "mount", "", "OIDC auth mount path (default: oidc)")
	oidcCmd.Flags().StringVar(&oidcRole, "role", "", "OIDC role name")
	oidcCmd.Flags().BoolVar(&oidcExport, "export", false, "Print token as export statement")

	rootCmd.AddCommand(oidcCmd)
}

func runOIDCLogin(cmd *cobra.Command, args []string) error {
	jwt := os.Getenv("OIDC_JWT_TOKEN")
	if jwt == "" {
		return fmt.Errorf("OIDC_JWT_TOKEN environment variable is required")
	}

	client, err := vault.NewClient("", "")
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	resp, err := vault.OIDCLogin(client, jwt, oidcRole, oidcMount)
	if err != nil {
		return fmt.Errorf("oidc login: %w", err)
	}

	if oidcExport {
		fmt.Fprintf(cmd.OutOrStdout(), "export VAULT_TOKEN=%s\n", resp.Token)
	} else {
		fmt.Fprintln(cmd.OutOrStdout(), resp.Token)
	}
	return nil
}
