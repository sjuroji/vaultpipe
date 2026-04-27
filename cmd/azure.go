package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/vault"
)

var (
	azureRole  string
	azureJWT   string
	azureMount string
	azureExport bool
)

var azureCmd = &cobra.Command{
	Use:   "azure",
	Short: "Authenticate with Vault using Azure auth method",
	RunE:  runAzureLogin,
}

func init() {
	azureCmd.Flags().StringVar(&azureRole, "role", "", "Azure role name (required)")
	azureCmd.Flags().StringVar(&azureJWT, "jwt", "", "Azure JWT token (required)")
	azureCmd.Flags().StringVar(&azureMount, "mount", "azure", "Auth mount path")
	azureCmd.Flags().BoolVar(&azureExport, "export", false, "Print token as export statement")
	_ = azureCmd.MarkFlagRequired("role")
	_ = azureCmd.MarkFlagRequired("jwt")
	rootCmd.AddCommand(azureCmd)
}

func runAzureLogin(cmd *cobra.Command, args []string) error {
	client, err := vault.NewClient("", "")
	if err != nil {
		return fmt.Errorf("azure login: %w", err)
	}

	token, err := vault.AzureLogin(client, azureRole, azureJWT, azureMount)
	if err != nil {
		return fmt.Errorf("azure login: %w", err)
	}

	if azureExport {
		fmt.Fprintf(os.Stdout, "export VAULT_TOKEN=%s\n", token)
	} else {
		fmt.Fprintln(os.Stdout, token)
	}
	return nil
}
