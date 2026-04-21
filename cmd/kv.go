package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/vault"
)

var kvCmd = &cobra.Command{
	Use:   "kv",
	Short: "Interact with Vault KV secrets engine",
	Long:  `List or delete secrets stored in a Vault KV v2 secrets engine.`,
}

var kvListCmd = &cobra.Command{
	Use:   "list <path>",
	Short: "List secrets at a KV path",
	Args:  cobra.ExactArgs(1),
	RunE:  runKVList,
}

var kvDeleteCmd = &cobra.Command{
	Use:   "delete <path> <key>",
	Short: "Delete a secret at a KV path",
	Args:  cobra.ExactArgs(2),
	RunE:  runKVDelete,
}

func init() {
	kvCmd.AddCommand(kvListCmd)
	kvCmd.AddCommand(kvDeleteCmd)
	RootCmd.AddCommand(kvCmd)
}

// runKVList lists all secret keys under the given KV path.
func runKVList(cmd *cobra.Command, args []string) error {
	addr, _ := cmd.Flags().GetString("vault-addr")
	token, _ := cmd.Flags().GetString("vault-token")

	if addr == "" {
		addr = os.Getenv("VAULT_ADDR")
	}
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}

	client, err := vault.NewClient(addr, token)
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	path := args[0]
	keys, err := vault.ListKV(client, path)
	if err != nil {
		return fmt.Errorf("failed to list KV path %q: %w", path, err)
	}

	if len(keys) == 0 {
		fmt.Fprintln(cmd.OutOrStdout(), "(no keys found)")
		return nil
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Keys at %s:\n", path)
	for _, k := range keys {
		fmt.Fprintf(cmd.OutOrStdout(), "  %s\n", strings.TrimSuffix(k, "/"))
	}
	return nil
}

// runKVDelete deletes the secret at the given KV path and key.
func runKVDelete(cmd *cobra.Command, args []string) error {
	addr, _ := cmd.Flags().GetString("vault-addr")
	token, _ := cmd.Flags().GetString("vault-token")

	if addr == "" {
		addr = os.Getenv("VAULT_ADDR")
	}
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}

	client, err := vault.NewClient(addr, token)
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	path := args[0]
	key := args[1]

	if err := vault.DeleteKV(client, path, key); err != nil {
		return fmt.Errorf("failed to delete KV secret %q at %q: %w", key, path, err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Deleted secret %q at path %q\n", key, path)
	return nil
}
