package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/vault"
)

func init() {
	policyCmd := &cobra.Command{
		Use:   "policy",
		Short: "Manage Vault ACL policies",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List all ACL policies",
		RunE:  runPolicyList,
	}

	getCmd := &cobra.Command{
		Use:   "get <name>",
		Short: "Get a specific ACL policy",
		Args:  cobra.ExactArgs(1),
		RunE:  runPolicyGet,
	}

	policyCmd.AddCommand(listCmd, getCmd)
	rootCmd.AddCommand(policyCmd)
}

func runPolicyList(cmd *cobra.Command, args []string) error {
	addr := os.Getenv("VAULT_ADDR")
	token := os.Getenv("VAULT_TOKEN")

	client, err := vault.NewClient(addr, token)
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}

	result, err := vault.ListPolicies(client)
	if err != nil {
		return fmt.Errorf("list policies: %w", err)
	}

	fmt.Println(strings.Join(result.Policies, "\n"))
	return nil
}

func runPolicyGet(cmd *cobra.Command, args []string) error {
	addr := os.Getenv("VAULT_ADDR")
	token := os.Getenv("VAULT_TOKEN")

	client, err := vault.NewClient(addr, token)
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}

	policy, err := vault.GetPolicy(client, args[0])
	if err != nil {
		return fmt.Errorf("get policy: %w", err)
	}

	fmt.Printf("Name:  %s\nRules:\n%s\n", policy.Name, policy.Rules)
	return nil
}
