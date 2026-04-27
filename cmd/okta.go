package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/vault"
)

var (
	oktaMount    string
	oktaUsername string
)

var oktaCmd = &cobra.Command{
	Use:   "okta",
	Short: "Authenticate with Vault using Okta credentials",
	RunE:  runOktaLogin,
}

func init() {
	oktaCmd.Flags().StringVar(&oktaMount, "mount", "okta", "Okta auth mount path")
	oktaCmd.Flags().StringVar(&oktaUsername, "username", "", "Okta username (required)")
	_ = oktaCmd.MarkFlagRequired("username")
	rootCmd.AddCommand(oktaCmd)
}

func runOktaLogin(cmd *cobra.Command, args []string) error {
	password := os.Getenv("OKTA_PASSWORD")
	if password == "" {
		return fmt.Errorf("OKTA_PASSWORD environment variable is not set")
	}

	client, err := vault.NewClient("", "")
	if err != nil {
		return fmt.Errorf("okta login: init client: %w", err)
	}

	token, err := vault.OktaLogin(client, oktaUsername, password, oktaMount)
	if err != nil {
		return fmt.Errorf("okta login: %w", err)
	}

	export, _ := cmd.Flags().GetBool("export")
	if export {
		fmt.Printf("export VAULT_TOKEN=%s\n", token)
	} else {
		fmt.Println(token)
	}

	return nil
}
