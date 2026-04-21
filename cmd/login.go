package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpipe/internal/vault"
)

var (
	loginUsername  string
	loginPassword  string
	loginMount     string
	loginExportEnv bool
)

func init() {
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with Vault using the userpass method",
		RunE:  runLogin,
	}

	loginCmd.Flags().StringVarP(&loginUsername, "username", "u", "", "Vault username (required)")
	loginCmd.Flags().StringVarP(&loginPassword, "password", "p", "", "Vault password (required)")
	loginCmd.Flags().StringVar(&loginMount, "mount", "userpass", "Auth mount path")
	loginCmd.Flags().BoolVar(&loginExportEnv, "export", false, "Print token as VAULT_TOKEN export statement")

	_ = loginCmd.MarkFlagRequired("username")
	_ = loginCmd.MarkFlagRequired("password")

	rootCmd.AddCommand(loginCmd)
}

func runLogin(cmd *cobra.Command, args []string) error {
	client, err := vault.NewClient("", "")
	if err != nil {
		return fmt.Errorf("login: init client: %w", err)
	}

	resp, err := vault.UserpassLogin(client, loginUsername, loginPassword, loginMount)
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}

	if loginExportEnv {
		fmt.Fprintf(os.Stdout, "export VAULT_TOKEN=%s\n", resp.Token)
	} else {
		fmt.Fprintf(os.Stdout, "Token:     %s\n", resp.Token)
		fmt.Fprintf(os.Stdout, "Renewable: %v\n", resp.Renewable)
		fmt.Fprintf(os.Stdout, "Lease:     %ds\n", resp.LeaseDur)
	}
	return nil
}
