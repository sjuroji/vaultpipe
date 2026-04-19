package runner

import (
	"os"

	"github.com/your-org/vaultpipe/internal/vault"
)

// EnvSource holds all resolved secrets and env-file variables that should be
// injected into the subprocess.
type EnvSource struct {
	// VaultSecrets maps environment variable names to Vault secret values.
	VaultSecrets map[string]string
	// EnvFileVars contains variables parsed from an .env file.
	EnvFileVars map[string]string
	// InheritOS controls whether the current process environment is inherited.
	InheritOS bool
}

// Build constructs the final []string environment slice for exec.
// Precedence (highest to lowest): VaultSecrets > EnvFileVars > OS environment.
func Build(src EnvSource) []string {
	base := make(map[string]string)

	if src.InheritOS {
		for _, pair := range os.Environ() {
			k, v := splitPair(pair)
			base[k] = v
		}
	}

	vault.MergeIntoEnv(base, src.EnvFileVars)
	vault.MergeIntoEnv(base, src.VaultSecrets)

	return vault.ToSlice(base)
}

// splitPair splits a "KEY=VALUE" string into its components.
func splitPair(pair string) (string, string) {
	for i := 0; i < len(pair); i++ {
		if pair[i] == '=' {
			return pair[:i], pair[i+1:]
		}
	}
	return pair, ""
}
