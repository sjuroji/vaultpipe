package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/vault"
)

var (
	certMount string
	certRole  string
	certExport bool
)

var certCmd = &cobra.Command{
	Use:   "cert",
	Short: "Authenticate to Vault using a TLS client certificate",
	RunE:  runCertLogin,
}

func init() {
	certCmd.Flags().StringVar(&certMount, "mount", "", "Auth mount path (default: cert)")
	certCmd.Flags().StringVar(&certRole, "role", "", "Certificate role name (optional)")
	certCmd.Flags().BoolVar(&certExport, "export", false, "Print token as export statement")
	rootCmd.AddCommand(certCmd)
}

func runCertLogin(cmd *cobra.Command, args []string) error {
	client, err := vault.NewClient(vaultAddr, vaultToken)
	if err != nil {
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	resp, err := vault.CertLogin(client, certMount, certRole)
	if err != nil {
		return fmt.Errorf("cert login failed: %w", err)
	}

	if certExport {
		fmt.Fprintf(os.Stdout, "export VAULT_TOKEN=%s\n", resp.Token)
	} else {
		fmt.Fprintf(os.Stdout, "token: %s\n", resp.Token)
		fmt.Fprintf(os.Stdout, "lease_duration: %d\n", resp.LeaseDur)
		fmt.Fprintf(os.Stdout, "renewable: %v\n", resp.Renewable)
	}
	return nil
}
