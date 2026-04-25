package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/vault"
)

var (
	approleMount  string
	approleExport bool
)

var approleCmd = &cobra.Command{
	Use:   "approle",
	Short: "Authenticate with Vault using AppRole",
	RunE:  runAppRoleLogin,
}

func init() {
	approleCmd.Flags().StringVar(&approleMount, "mount", "approle", "AppRole auth mount path")
	approleCmd.Flags().BoolVar(&approleExport, "export", false, "Print token as export statement")
	rootCmd.AddCommand(approleCmd)
}

func runAppRoleLogin(cmd *cobra.Command, args []string) error {
	roleID := os.Getenv("VAULT_ROLE_ID")
	if roleID == "" {
		return fmt.Errorf("VAULT_ROLE_ID environment variable is not set")
	}

	secretID := os.Getenv("VAULT_SECRET_ID")
	if secretID == "" {
		return fmt.Errorf("VAULT_SECRET_ID environment variable is not set")
	}

	addr, _ := cmd.Root().PersistentFlags().GetString("address")
	token, _ := cmd.Root().PersistentFlags().GetString("token")

	c, err := vault.NewClient(addr, token)
	if err != nil {
		return fmt.Errorf("approle: create client: %w", err)
	}

	clientToken, err := vault.AppRoleLogin(c, roleID, secretID, approleMount)
	if err != nil {
		return fmt.Errorf("approle login failed: %w", err)
	}

	if approleExport {
		fmt.Fprintf(cmd.OutOrStdout(), "export VAULT_TOKEN=%s\n", clientToken)
	} else {
		fmt.Fprintln(cmd.OutOrStdout(), clientToken)
	}

	return nil
}
