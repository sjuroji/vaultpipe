package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/vault"
)

var (
	jwtMount string
	jwtRole  string
	jwtToken string
	jwtExport bool
)

func init() {
	jwtCmd := &cobra.Command{
		Use:   "jwt",
		Short: "Authenticate with Vault using a JWT/OIDC token",
		RunE:  runJWTLogin,
	}

	jwtCmd.Flags().StringVar(&jwtMount, "mount", "", "Auth mount path (default: jwt)")
	jwtCmd.Flags().StringVar(&jwtRole, "role", "", "Vault role to authenticate against (required)")
	jwtCmd.Flags().StringVar(&jwtToken, "jwt", "", "JWT token (overrides VAULT_JWT env var)")
	jwtCmd.Flags().BoolVar(&jwtExport, "export", false, "Print export statement instead of plain token")
	_ = jwtCmd.MarkFlagRequired("role")

	rootCmd.AddCommand(jwtCmd)
}

func runJWTLogin(cmd *cobra.Command, args []string) error {
	token := jwtToken
	if token == "" {
		token = os.Getenv("VAULT_JWT")
	}
	if token == "" {
		return fmt.Errorf("jwt token required: set --jwt or VAULT_JWT")
	}

	client, err := vault.NewClient("", "")
	if err != nil {
		return fmt.Errorf("create vault client: %w", err)
	}

	resp, err := vault.JWTLogin(client, token, jwtRole, jwtMount)
	if err != nil {
		return fmt.Errorf("jwt login: %w", err)
	}

	if jwtExport {
		fmt.Printf("export VAULT_TOKEN=%s\n", resp.Token)
	} else {
		fmt.Println(resp.Token)
	}
	return nil
}
