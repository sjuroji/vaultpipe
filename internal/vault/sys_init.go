package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// InitStatus holds the initialization state of the Vault cluster.
type InitStatus struct {
	Initialized bool `json:"initialized"`
}

// InitRequest holds the parameters for initializing a Vault cluster.
type InitRequest struct {
	SecretShares    int `json:"secret_shares"`
	SecretThreshold int `json:"secret_threshold"`
}

// InitResponse holds the root token and unseal keys returned on init.
type InitResponse struct {
	RootToken  string   `json:"root_token"`
	UnsealKeys []string `json:"keys"`
}

// GetInitStatus returns whether the Vault instance has been initialized.
func GetInitStatus(client *Client) (*InitStatus, error) {
	resp, err := client.Get("/v1/sys/init")
	if err != nil {
		return nil, fmt.Errorf("init status request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var status InitStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode init status: %w", err)
	}
	return &status, nil
}

// Initialize initializes a Vault cluster with the given shares and threshold.
func Initialize(client *Client, shares, threshold int) (*InitResponse, error) {
	body := InitRequest{
		SecretShares:    shares,
		SecretThreshold: threshold,
	}

	resp, err := client.Post("/v1/sys/init", body)
	if err != nil {
		return nil, fmt.Errorf("init request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var initResp InitResponse
	if err := json.NewDecoder(resp.Body).Decode(&initResp); err != nil {
		return nil, fmt.Errorf("failed to decode init response: %w", err)
	}
	return &initResp, nil
}
