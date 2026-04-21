package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yourorg/vaultpipe/internal/vault"
)

var (
	rekeyShares    int
	rekeyThreshold int
)

var rekeyCmd = &cobra.Command{
	Use:   "rekey",
	Short: "Manage Vault rekey operations",
}

var rekeyInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Start a new rekey operation",
	RunE:  runRekeyInit,
}

var rekeyStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Get current rekey operation status",
	RunE:  runRekeyStatus,
}

var rekeyCancelCmd = &cobra.Command{
	Use:   "cancel",
	Short: "Cancel an in-progress rekey operation",
	RunE:  runRekeyCancel,
}

func init() {
	rekeyInitCmd.Flags().IntVar(&rekeyShares, "shares", 5, "Number of key shares")
	rekeyInitCmd.Flags().IntVar(&rekeyThreshold, "threshold", 3, "Number of key shares required to unseal")
	rekeyCmd.AddCommand(rekeyInitCmd, rekeyStatusCmd, rekeyCancelCmd)
	rootCmd.AddCommand(rekeyCmd)
}

func runRekeyInit(cmd *cobra.Command, args []string) error {
	c, err := vault.NewClient("", "")
	if err != nil {
		return err
	}
	status, err := vault.InitRekey(c, rekeyShares, rekeyThreshold)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "Rekey started: nonce=%s shares=%d threshold=%d\n",
		status.Nonce, status.N, status.T)
	return nil
}

func runRekeyStatus(cmd *cobra.Command, args []string) error {
	c, err := vault.NewClient("", "")
	if err != nil {
		return err
	}
	status, err := vault.GetRekeyStatus(c)
	if err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "started=%v nonce=%s progress=%d/%d\n",
		status.Started, status.Nonce, status.Progress, status.Required)
	return nil
}

func runRekeyCancel(cmd *cobra.Command, args []string) error {
	c, err := vault.NewClient("", "")
	if err != nil {
		return err
	}
	if err := vault.CancelRekey(c); err != nil {
		return err
	}
	fmt.Fprintln(os.Stdout, "Rekey operation cancelled")
	return nil
}
