package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	envFile    string
	vaultPath  string
	vaultAddr  string
	vaultToken string
	inheritEnv bool
)

var rootCmd = &cobra.Command{
	Use:   "vaultpipe [flags] -- <command> [args...]",
	Short: "Inject secrets from Vault or env files into subprocess environments",
	Args:  cobra.MinimumNArgs(1),
	RunE:  runRoot,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&envFile, "env-file", "e", "", "Path to .env file")
	rootCmd.Flags().StringVarP(&vaultPath, "vault-path", "p", "", "Vault KV secret path (e.g. secret/data/myapp)")
	rootCmd.Flags().StringVar(&vaultAddr, "vault-addr", "", "Vault server address (overrides VAULT_ADDR)")
	rootCmd.Flags().StringVar(&vaultToken, "vault-token", "", "Vault token (overrides VAULT_TOKEN)")
	rootCmd.Flags().BoolVar(&inheritEnv, "inherit-env", true, "Inherit host environment variables")
}
