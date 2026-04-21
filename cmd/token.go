package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpipe/internal/vault"
)

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Display information about the current Vault token",
	RunE:  runToken,
}

func init() {
	rootCmd.AddCommand(tokenCmd)
}

func runToken(cmd *cobra.Command, args []string) error {
	address, _ := cmd.Root().PersistentFlags().GetString("address")
	token, _ := cmd.Root().PersistentFlags().GetString("token")

	if address == "" {
		address = os.Getenv("VAULT_ADDR")
	}
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}

	client, err := vault.NewClient(address, token)
	if err != nil {
		return fmt.Errorf("create client: %w", err)
	}

	info, err := vault.LookupSelfToken(client)
	if err != nil {
		return fmt.Errorf("lookup token: %w", err)
	}

	w := tabwriter.NewWriter(cmd.OutOrStdout(), 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "ID:\t%s\n", info.ID)
	fmt.Fprintf(w, "Display Name:\t%s\n", info.DisplayName)
	fmt.Fprintf(w, "TTL:\t%s\n", info.TTL)
	fmt.Fprintf(w, "Renewable:\t%v\n", info.Renewable)
	fmt.Fprintf(w, "Policies:\t%v\n", info.Policies)
	fmt.Fprintf(w, "Created:\t%s\n", info.CreationTime.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(w, "Expired:\t%v\n", info.IsExpired())
	return w.Flush()
}
