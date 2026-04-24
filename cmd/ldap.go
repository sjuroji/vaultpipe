package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/vault"
)

var (
	ldapMount    string
	ldapUsername string
	ldapExport   bool
)

var ldapCmd = &cobra.Command{
	Use:   "ldap",
	Short: "Authenticate with Vault using LDAP credentials",
	RunE:  runLDAPLogin,
}

func init() {
	ldapCmd.Flags().StringVar(&ldapMount, "mount", "ldap", "LDAP auth mount path")
	ldapCmd.Flags().StringVar(&ldapUsername, "username", "", "LDAP username (required)")
	ldapCmd.Flags().BoolVar(&ldapExport, "export", false, "Print token as export statement")
	_ = ldapCmd.MarkFlagRequired("username")
	rootCmd.AddCommand(ldapCmd)
}

func runLDAPLogin(cmd *cobra.Command, args []string) error {
	password := os.Getenv("LDAP_PASSWORD")
	if password == "" {
		return fmt.Errorf("LDAP_PASSWORD environment variable is not set")
	}

	client, err := vault.NewClient(vaultAddr, vaultToken)
	if err != nil {
		return fmt.Errorf("ldap login: create client: %w", err)
	}

	resp, err := vault.LDAPLogin(client, ldapUsername, password, ldapMount)
	if err != nil {
		return fmt.Errorf("ldap login: %w", err)
	}

	if ldapExport {
		fmt.Fprintf(cmd.OutOrStdout(), "export VAULT_TOKEN=%s\n", resp.Token)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "token: %s\npolicies: %v\nlease_duration: %ds\n",
			resp.Token, resp.Policies, resp.LeaseDur)
	}
	return nil
}
