package vault

import "strings"

// FilterByPrefix returns a new map containing only keys that start with the given prefix.
func FilterByPrefix(secrets map[string]string, prefix string) map[string]string {
	if prefix == "" {
		return secrets
	}
	filtered := make(map[string]string)
	for k, v := range secrets {
		if strings.HasPrefix(k, prefix) {
			filtered[k] = v
		}
	}
	return filtered
}

// AddPrefix returns a new map with the given prefix prepended to every key.
func AddPrefix(secrets map[string]string, prefix string) map[string]string {
	if prefix == "" {
		return secrets
	}
	prefixed := make(map[string]string, len(secrets))
	for k, v := range secrets {
		prefixed[prefix+k] = v
	}
	return prefixed
}

// StripPrefix returns a new map with the given prefix removed from every key
// that starts with it. Keys that do not start with the prefix are left unchanged.
func StripPrefix(secrets map[string]string, prefix string) map[string]string {
	if prefix == "" {
		return secrets
	}
	stripped := make(map[string]string, len(secrets))
	for k, v := range secrets {
		stripped[strings.TrimPrefix(k, prefix)] = v
	}
	return stripped
}
