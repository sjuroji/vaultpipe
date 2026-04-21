package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpipe/internal/vault"
)

var mountCmd = &cobra.Command{
	Use:   "mount [path]",
	Short: "List or inspect Vault secrets engine mounts",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runMount,
}

func init() {
	rootCmd.AddCommand(mountCmd)
}

func runMount(cmd *cobra.Command, args []string) error {
	addr, _ := cmd.Flags().GetString("vault-addr")
	token, _ := cmd.Flags().GetString("vault-token")

	c, err := vault.NewClient(addr, token)
	if err != nil {
		return fmt.Errorf("init vault client: %w", err)
	}

	if len(args) == 1 {
		info, err := vault.GetMount(c, args[0])
		if err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "Path:        %s\n", args[0])
		fmt.Fprintf(cmd.OutOrStdout(), "Type:        %s\n", info.Type)
		fmt.Fprintf(cmd.OutOrStdout(), "Description: %s\n", info.Description)
		fmt.Fprintf(cmd.OutOrStdout(), "Accessor:    %s\n", info.Accessor)
		fmt.Fprintf(cmd.OutOrStdout(), "Local:       %v\n", info.Local)
		fmt.Fprintf(cmd.OutOrStdout(), "Seal Wrap:   %v\n", info.SealWrap)
		return nil
	}

	mounts, err := vault.ListMounts(c)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PATH\tTYPE\tDESCRIPTION\tACCESSOR")
	for path, info := range mounts {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", path, info.Type, info.Description, info.Accessor)
	}
	return w.Flush()
}
