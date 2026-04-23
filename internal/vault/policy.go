package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// CapabilitiesResponse holds the capabilities for a given path.
type CapabilitiesResponse struct {
	Capabilities []string `json:"capabilities"`
}

// CheckCapabilities queries Vault for the capabilities of the current token on a path.
func CheckCapabilities(client *Client, path string) (*CapabilitiesResponse, error) {
	payload := map[string]string{"path": path}
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	resp, err := client.PostJSON("/v1/sys/capabilities-self", body)
	if err != nil {
		return nil, fmt.Errorf("capabilities request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result CapabilitiesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}

// HasCapability returns true if the given capability is present in the response.
func HasCapability(caps *CapabilitiesResponse, capability string) bool {
	for _, c := range caps.Capabilities {
		if c == capability {
			return true
		}
	}
	return false
}
