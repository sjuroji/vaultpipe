package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/vault"
)

var (
	unsealReset bool
)

func init() {
	unsealCmd := &cobra.Command{
		Use:   "unseal [key]",
		Short: "Submit an unseal key shard to Vault",
		Long: `Submit a single unseal key shard to progress the unseal operation.
Omit [key] and use --reset to cancel an in-progress unseal attempt.`,
		RunE: runUnseal,
	}

	unsealCmd.Flags().BoolVar(&unsealReset, "reset", false, "Reset (cancel) an in-progress unseal attempt")
	rootCmd.AddCommand(unsealCmd)
}

func runUnseal(cmd *cobra.Command, args []string) error {
	key := ""
	if !unsealReset {
		if len(args) < 1 {
			return fmt.Errorf("unseal key required (or use --reset to cancel)")
		}
		key = args[0]
	}

	c, err := vault.NewClient(vaultAddr, vaultToken)
	if err != nil {
		return fmt.Errorf("unseal: create client: %w", err)
	}

	res, err := vault.SubmitUnsealKey(c, key, unsealReset)
	if err != nil {
		return fmt.Errorf("unseal: %w", err)
	}

	if unsealReset {
		fmt.Fprintln(os.Stdout, "Unseal attempt reset. Progress cleared.")
		return nil
	}

	if res.Sealed {
		fmt.Fprintf(os.Stdout, "Progress: %d/%d — Vault is still sealed.\n", res.Progress, res.Threshold)
	} else {
		fmt.Fprintln(os.Stdout, "Vault is now unsealed.")
	}
	return nil
}
