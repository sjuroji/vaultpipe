package vault

import "fmt"

// MergeIntoEnv merges secret key/value pairs into an existing environment map.
// Existing keys are preserved unless overwrite is true.
func MergeIntoEnv(env map[string]string, secrets map[string]string, overwrite bool) map[string]string {
	out := make(map[string]string, len(env)+len(secrets))
	for k, v := range env {
		out[k] = v
	}
	for k, v := range secrets {
		if _, exists := out[k]; !exists || overwrite {
			out[k] = v
		}
	}
	return out
}

// ToSlice converts an env map to a []string slice in "KEY=VALUE" format
// suitable for use with exec.Cmd.Env.
func ToSlice(env map[string]string) []string {
	slice := make([]string, 0, len(env))
	for k, v := range env {
		slice = append(slice, k+"="+v)
	}
	return slice
}

// FromSlice parses a []string slice of "KEY=VALUE" pairs into a map.
// Entries that do not contain "=" are skipped.
func FromSlice(env []string) map[string]string {
	out := make(map[string]string, len(env))
	for _, entry := range env {
		for i := 0; i < len(entry); i++ {
			if entry[i] == '=' {
				out[entry[:i]] = entry[i+1:]
				break
			}
		}
	}
	return out
}

// formatEntry formats a single key/value pair as "KEY=VALUE".
func formatEntry(k, v string) string {
	return fmt.Sprintf("%s=%s", k, v)
}
