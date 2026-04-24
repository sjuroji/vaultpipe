package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/vault"
)

var (
	githubMount string
	githubExport bool
)

func init() {
	githubCmd := &cobra.Command{
		Use:   "github",
		Short: "Authenticate with Vault using a GitHub personal access token",
		RunE:  runGitHubLogin,
	}

	githubCmd.Flags().StringVar(&githubMount, "mount", "github", "GitHub auth mount path")
	githubCmd.Flags().BoolVar(&githubExport, "export", false, "Print token as export statement")

	rootCmd.AddCommand(githubCmd)
}

func runGitHubLogin(cmd *cobra.Command, args []string) error {
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		return fmt.Errorf("GITHUB_TOKEN environment variable is not set")
	}

	c, err := vault.NewClient("", "")
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	resp, err := vault.GitHubLogin(c, githubToken, githubMount)
	if err != nil {
		return fmt.Errorf("github login: %w", err)
	}

	if githubExport {
		fmt.Fprintf(cmd.OutOrStdout(), "export VAULT_TOKEN=%s\n", resp.Token)
	} else {
		fmt.Fprintf(cmd.OutOrStdout(), "token: %s\n", resp.Token)
		fmt.Fprintf(cmd.OutOrStdout(), "lease_duration: %d\n", resp.LeaseDur)
		fmt.Fprintf(cmd.OutOrStdout(), "renewable: %v\n", resp.Renewable)
	}
	return nil
}
