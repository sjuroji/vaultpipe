package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/vault"
)

var (
	radiusMount    string
	radiusUsername string
	radiusPassword string
	radiusExport   bool
)

func init() {
	radiusCmd := &cobra.Command{
		Use:   "radius",
		Short: "Authenticate using the RADIUS auth method",
		RunE:  runRADIUSLogin,
	}
	radiusCmd.Flags().StringVar(&radiusMount, "mount", "radius", "RADIUS auth mount path")
	radiusCmd.Flags().StringVar(&radiusUsername, "username", "", "RADIUS username")
	radiusCmd.Flags().StringVar(&radiusPassword, "password", "", "RADIUS password (or set RADIUS_PASSWORD)")
	radiusCmd.Flags().BoolVar(&radiusExport, "export", false, "Print token as export statement")
	_ = radiusCmd.MarkFlagRequired("username")
	rootCmd.AddCommand(radiusCmd)
}

func runRADIUSLogin(cmd *cobra.Command, args []string) error {
	pwd := radiusPassword
	if pwd == "" {
		pwd = os.Getenv("RADIUS_PASSWORD")
	}
	if pwd == "" {
		return fmt.Errorf("radius login: password required via --password or RADIUS_PASSWORD")
	}

	c, err := vault.NewClient("", "")
	if err != nil {
		return fmt.Errorf("radius login: %w", err)
	}

	token, err := vault.RADIUSLogin(c, radiusUsername, pwd, radiusMount)
	if err != nil {
		return err
	}

	if radiusExport {
		fmt.Fprintf(cmd.OutOrStdout(), "export VAULT_TOKEN=%s\n", token)
	} else {
		fmt.Fprintln(cmd.OutOrStdout(), token)
	}
	return nil
}
