package runner

import (
	"os"
	"strings"

	"github.com/your-org/vaultpipe/internal/envfile"
	"github.com/your-org/vaultpipe/internal/vault"
)

// Options controls how the subprocess environment is assembled.
type Options struct {
	InheritOS   bool
	EnvFilePath string
	VaultSecrets map[string]string // key -> vault path
	VaultClient  *vault.Client
}

// Build assembles the final environment slice for the subprocess.
// Priority (highest to lowest): Vault secrets > env file > OS environment.
func Build(opts Options) ([]string, error) {
	env := map[string]string{}

	if opts.InheritOS {
		for _, pair := range os.Environ() {
			k, v := splitPair(pair)
			env[k] = v
		}
	}

	if opts.EnvFilePath != "" {
		parsed, err := envfile.Parse(opts.EnvFilePath)
		if err != nil {
			return nil, err
		}
		for k, v := range parsed {
			env[k] = v
		}
	}

	if opts.VaultClient != nil {
		for envKey, secretPath := range opts.VaultSecrets {
			secrets, err := opts.VaultClient.ReadSecret(secretPath)
			if err != nil {
				return nil, err
			}
			vault.MergeIntoEnv(env, secrets)
			_ = envKey // path-level merge; envKey reserved for future field mapping
		}
	}

	return vault.ToSlice(env), nil
}

func splitPair(pair string) (string, string) {
	parts := strings.SplitN(pair, "=", 2)
	if len(parts) != 2 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}
