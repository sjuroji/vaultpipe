package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/yourorg/vaultpipe/internal/vault"
)

var snapshotOutput string

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Manage Vault raft snapshots",
}

var snapshotTakeCmd = &cobra.Command{
	Use:   "take",
	Short: "Take a raft snapshot and write it to a file",
	RunE:  runSnapshotTake,
}

var snapshotStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current raft snapshot metadata",
	RunE:  runSnapshotStatus,
}

func init() {
	snapshotTakeCmd.Flags().StringVarP(&snapshotOutput, "output", "o", fmt.Sprintf("vault-snapshot-%d.snap", time.Now().Unix()), "Output file path for the snapshot")
	snapshotCmd.AddCommand(snapshotTakeCmd)
	snapshotCmd.AddCommand(snapshotStatusCmd)
	rootCmd.AddCommand(snapshotCmd)
}

func runSnapshotTake(cmd *cobra.Command, args []string) error {
	c, err := vault.NewClient("", "")
	if err != nil {
		return fmt.Errorf("snapshot take: %w", err)
	}

	f, err := os.Create(snapshotOutput)
	if err != nil {
		return fmt.Errorf("snapshot take: create file: %w", err)
	}
	defer f.Close()

	n, err := vault.TakeSnapshot(c, f)
	if err != nil {
		return fmt.Errorf("snapshot take: %w", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Snapshot written to %s (%d bytes)\n", snapshotOutput, n)
	return nil
}

func runSnapshotStatus(cmd *cobra.Command, args []string) error {
	c, err := vault.NewClient("", "")
	if err != nil {
		return fmt.Errorf("snapshot status: %w", err)
	}

	meta, err := vault.SnapshotStatus(c)
	if err != nil {
		return fmt.Errorf("snapshot status: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Index:     %d\nTerm:      %d\nVersion:   %d\nTimestamp: %s\n",
		meta.Index, meta.Term, meta.Version, meta.Timestamp.Format(time.RFC3339))
	return nil
}
