package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/example/vaultpipe/internal/runner"
	"github.com/example/vaultpipe/internal/vault"
)

func runRoot(cmd *cobra.Command, args []string) error {
	var vaultSecrets map[string]string

	if vaultPath != "" {
		client, err := vault.NewClient(vaultAddr, vaultToken)
		if err != nil {
			return fmt.Errorf("vault client: %w", err)
		}
		vaultSecrets, err = client.ReadSecret(vaultPath)
		if err != nil {
			return fmt.Errorf("vault read: %w", err)
		}
	}

	env, err := runner.Build(runner.BuildConfig{
		InheritOS:    inheritEnv,
		EnvFilePath:  envFile,
		VaultSecrets: vaultSecrets,
	})
	if err != nil {
		return fmt.Errorf("build env: %w", err)
	}

	code, err := runner.Run(args, env)
	if err != nil {
		return fmt.Errorf("exec: %w", err)
	}
	if code != 0 {
		os.Exit(code)
	}
	return nil
}
