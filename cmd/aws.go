package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/vault"
)

var (
	awsRole            string
	awsMount           string
	awsIAMRequestURL   string
	awsIAMRequestBody  string
	awsIAMHeaders      string
	awsExport          bool
)

func init() {
	awsCmd := &cobra.Command{
		Use:   "aws",
		Short: "Authenticate with Vault using the AWS IAM auth method",
		RunE:  runAWSLogin,
	}

	awsCmd.Flags().StringVar(&awsRole, "role", "", "AWS IAM role name configured in Vault (required)")
	awsCmd.Flags().StringVar(&awsMount, "mount", "aws", "Vault AWS auth mount path")
	awsCmd.Flags().StringVar(&awsIAMRequestURL, "iam-request-url", "", "Base64-encoded IAM request URL")
	awsCmd.Flags().StringVar(&awsIAMRequestBody, "iam-request-body", "", "Base64-encoded IAM request body")
	awsCmd.Flags().StringVar(&awsIAMHeaders, "iam-request-headers", "", "Base64-encoded IAM request headers")
	awsCmd.Flags().BoolVar(&awsExport, "export", false, "Print token as export statement")

	_ = awsCmd.MarkFlagRequired("role")

	rootCmd.AddCommand(awsCmd)
}

func runAWSLogin(cmd *cobra.Command, args []string) error {
	c, err := vault.NewClient("", "")
	if err != nil {
		return fmt.Errorf("aws login: create client: %w", err)
	}

	token, err := vault.AWSLogin(c, awsRole, awsIAMRequestURL, awsIAMRequestBody, awsIAMHeaders, awsMount)
	if err != nil {
		return fmt.Errorf("aws login: %w", err)
	}

	if awsExport {
		fmt.Fprintf(os.Stdout, "export VAULT_TOKEN=%s\n", token)
	} else {
		fmt.Fprintln(os.Stdout, token)
	}
	return nil
}
