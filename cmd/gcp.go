package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/vault"
)

var (
	gcpRole  string
	gcpJWT   string
	gcpMount string
	gcpExport bool
)

var gcpCmd = &cobra.Command{
	Use:   "gcp",
	Short: "Authenticate to Vault using GCP IAM",
	RunE:  runGCPLogin,
}

func init() {
	gcpCmd.Flags().StringVar(&gcpRole, "role", "", "GCP role name configured in Vault (required)")
	gcpCmd.Flags().StringVar(&gcpJWT, "jwt", "", "Signed GCP IAM JWT (required)")
	gcpCmd.Flags().StringVar(&gcpMount, "mount", "", "Auth mount path (default: gcp)")
	gcpCmd.Flags().BoolVar(&gcpExport, "export", false, "Print token as export statement")
	_ = gcpCmd.MarkFlagRequired("role")
	_ = gcpCmd.MarkFlagRequired("jwt")
	rootCmd.AddCommand(gcpCmd)
}

func runGCPLogin(cmd *cobra.Command, args []string) error {
	client, err := vault.NewClient("", "")
	if err != nil {
		return fmt.Errorf("create vault client: %w", err)
	}

	resp, err := vault.GCPLogin(client, vault.GCPLoginRequest{
		Role:  gcpRole,
		JWT:   gcpJWT,
		Mount: gcpMount,
	})
	if err != nil {
		return fmt.Errorf("gcp login: %w", err)
	}

	if gcpExport {
		fmt.Fprintf(os.Stdout, "export VAULT_TOKEN=%s\n", resp.ClientToken)
	} else {
		fmt.Fprintln(os.Stdout, resp.ClientToken)
	}
	return nil
}
