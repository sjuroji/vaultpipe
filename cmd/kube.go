package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/vault"
)

var (
	kubeRole   string
	kubeMount  string
	kubeExport bool
)

// kubeCmd is the parent command for Kubernetes auth operations.
var kubeCmd = &cobra.Command{
	Use:   "kube",
	Short: "Authenticate with Vault using Kubernetes service account",
	Long: `Login to Vault using the Kubernetes auth method.

The Kubernetes JWT token is read from the standard service account
token path (/var/run/secrets/kubernetes.io/serviceaccount/token) unless
overridden via the KUBE_TOKEN environment variable.`,
	RunE: runKubeLogin,
}

func init() {
	kubeCmd.Flags().StringVar(&kubeRole, "role", "", "Vault role to authenticate against (required)")
	kubeCmd.Flags().StringVar(&kubeMount, "mount", "kubernetes", "Vault mount path for the Kubernetes auth method")
	kubeCmd.Flags().BoolVar(&kubeExport, "export", false, "Print token as export statement instead of plain value")
	_ = kubeCmd.MarkFlagRequired("role")

	rootCmd.AddCommand(kubeCmd)
}

// runKubeLogin performs a Kubernetes auth login and prints the resulting token.
func runKubeLogin(cmd *cobra.Command, args []string) error {
	// Resolve the JWT token: prefer env override, fall back to the standard
	// Kubernetes service account token file.
	token := os.Getenv("KUBE_TOKEN")
	if token == "" {
		const saTokenPath = "/var/run/secrets/kubernetes.io/serviceaccount/token"
		data, err := os.ReadFile(saTokenPath)
		if err != nil {
			return fmt.Errorf("reading Kubernetes service account token from %s: %w", saTokenPath, err)
		}
		token = string(data)
	}

	client, err := vault.NewClient(vaultAddr, vaultToken)
	if err != nil {
		return fmt.Errorf("creating Vault client: %w", err)
	}

	vaultTok, err := vault.KubernetesLogin(client, kubeRole, token, kubeMount)
	if err != nil {
		return fmt.Errorf("Kubernetes login failed: %w", err)
	}

	if kubeExport {
		fmt.Fprintf(cmd.OutOrStdout(), "export VAULT_TOKEN=%s\n", vaultTok)
	} else {
		fmt.Fprintln(cmd.OutOrStdout(), vaultTok)
	}

	return nil
}
