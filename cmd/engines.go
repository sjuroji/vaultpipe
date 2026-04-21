package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpipe/internal/vault"
)

var enginesCmd = &cobra.Command{
	Use:   "engines",
	Short: "List or inspect Vault secret engines",
}

var enginesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all mounted secret engines",
	RunE:  runEnginesList,
}

var enginesGetCmd = &cobra.Command{
	Use:   "get <path>",
	Short: "Get details of a specific secret engine by mount path",
	Args:  cobra.ExactArgs(1),
	RunE:  runEnginesGet,
}

func init() {
	enginesCmd.AddCommand(enginesListCmd)
	enginesCmd.AddCommand(enginesGetCmd)
	rootCmd.AddCommand(enginesCmd)
}

func runEnginesList(cmd *cobra.Command, args []string) error {
	c, err := vault.NewClient("", "")
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	engines, err := vault.ListSecretEngines(c)
	if err != nil {
		return fmt.Errorf("list engines: %w", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "PATH\tTYPE\tDESCRIPTION\tACCESSOR")
	for path, eng := range engines {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", path, eng.Type, eng.Description, eng.Accessor)
	}
	return w.Flush()
}

func runEnginesGet(cmd *cobra.Command, args []string) error {
	c, err := vault.NewClient("", "")
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	engine, err := vault.GetSecretEngine(c, args[0])
	if err != nil {
		return fmt.Errorf("get engine: %w", err)
	}

	fmt.Printf("Type:        %s\n", engine.Type)
	fmt.Printf("Description: %s\n", engine.Description)
	fmt.Printf("Accessor:    %s\n", engine.Accessor)
	fmt.Printf("Local:       %v\n", engine.Local)
	fmt.Printf("Seal Wrap:   %v\n", engine.SealWrap)
	return nil
}
