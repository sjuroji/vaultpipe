package vault

import "strings"

// FilterByPrefix returns a new map containing only entries whose keys
// start with the given prefix. The prefix is stripped from the keys
// in the returned map.
func FilterByPrefix(env map[string]string, prefix string) map[string]string {
	out := make(map[string]string)
	for k, v := range env {
		if strings.HasPrefix(k, prefix) {
			newKey := strings.TrimPrefix(k, prefix)
			if newKey != "" {
				out[newKey] = v
			}
		}
	}
	return out
}

// AddPrefix returns a new map with the given prefix prepended to every key.
func AddPrefix(env map[string]string, prefix string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[prefix+k] = v
	}
	return out
}
