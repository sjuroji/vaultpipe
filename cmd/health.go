package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpipe/internal/vault"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check the health of the configured Vault server",
	Long: `Queries the Vault /v1/sys/health endpoint and prints the cluster status.
Exits with code 1 if Vault is sealed, uninitialized, or unreachable.`,
	RunE: runHealth,
}

func init() {
	rootCmd.AddCommand(healthCmd)
}

func runHealth(cmd *cobra.Command, _ []string) error {
	addr := os.Getenv("VAULT_ADDR")
	if addr == "" {
		addr = "http://127.0.0.1:8200"
	}

	client, err := vault.NewClient(addr, "")
	if err != nil {
		return fmt.Errorf("health: create client: %w", err)
	}

	status, err := vault.CheckHealth(context.Background(), client.HTTP(), addr)
	if err != nil {
		return fmt.Errorf("health: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Vault %s\n", status.Version)
	fmt.Fprintf(cmd.OutOrStdout(), "  initialized : %v\n", status.Initialized)
	fmt.Fprintf(cmd.OutOrStdout(), "  sealed      : %v\n", status.Sealed)
	fmt.Fprintf(cmd.OutOrStdout(), "  standby     : %v\n", status.Standby)

	if status.ClusterName != "" {
		fmt.Fprintf(cmd.OutOrStdout(), "  cluster     : %s\n", status.ClusterName)
	}

	if !vault.IsHealthy(status) {
		return fmt.Errorf("vault is not healthy (sealed=%v, standby=%v, initialized=%v)",
			status.Sealed, status.Standby, status.Initialized)
	}

	fmt.Fprintln(cmd.OutOrStdout(), "Status: OK")
	return nil
}
