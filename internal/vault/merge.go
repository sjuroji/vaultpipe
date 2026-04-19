package vault

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
