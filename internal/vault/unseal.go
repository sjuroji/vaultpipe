package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// UnsealRequest is the payload sent to the unseal endpoint.
type UnsealRequest struct {
	Key   string `json:"key"`
	Reset bool   `json:"reset,omitempty"`
}

// UnsealResponse holds the result of an unseal attempt.
type UnsealResponse struct {
	Sealed    bool `json:"sealed"`
	T         int  `json:"t"`
	N         int  `json:"n"`
	Progress  int  `json:"progress"`
	Threshold int  `json:"threshold"`
}

// SubmitUnsealKey submits a single unseal key shard to Vault.
// Returns the current unseal progress. If Reset is true, the unseal
// attempt is cancelled and progress is reset to zero.
func SubmitUnsealKey(c *Client, key string, reset bool) (*UnsealResponse, error) {
	body, err := json.Marshal(UnsealRequest{Key: key, Reset: reset})
	if err != nil {
		return nil, fmt.Errorf("unseal: marshal request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, c.Address+"/v1/sys/unseal", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("unseal: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unseal: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unseal: unexpected status %d", resp.StatusCode)
	}

	var result UnsealResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("unseal: decode response: %w", err)
	}
	return &result, nil
}

// ResetUnseal cancels any in-progress unseal attempt, resetting the key
// shard counter back to zero. It is a convenience wrapper around
// SubmitUnsealKey with Reset set to true and an empty key.
func ResetUnseal(c *Client) (*UnsealResponse, error) {
	return SubmitUnsealKey(c, "", true)
}
