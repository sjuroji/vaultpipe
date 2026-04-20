package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// PolicyCapabilities holds the capabilities a token has on a given path.
type PolicyCapabilities struct {
	Capabilities []string `json:"capabilities"`
}

// CheckCapabilities queries Vault to determine what capabilities the current
// token has on the given path. It returns a slice of capability strings such
// as "read", "list", "deny", etc.
func CheckCapabilities(c *Client, path string) ([]string, error) {
	body, err := json.Marshal(map[string]string{"path": path})
	if err != nil {
		return nil, fmt.Errorf("policy: marshal request: %w", err)
	}

	resp, err := c.post("/v1/sys/capabilities-self", body)
	if err != nil {
		return nil, fmt.Errorf("policy: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("policy: unexpected status %d", resp.StatusCode)
	}

	var result PolicyCapabilities
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("policy: decode response: %w", err)
	}

	return result.Capabilities, nil
}

// HasCapability returns true if the given capability (e.g. "read") is present
// in the list returned by CheckCapabilities for the specified path.
func HasCapability(c *Client, path, capability string) (bool, error) {
	caps, err := CheckCapabilities(c, path)
	if err != nil {
		return false, err
	}
	for _, cap := range caps {
		if cap == capability {
			return true, nil
		}
	}
	return false, nil
}
