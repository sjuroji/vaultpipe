package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SealStatus represents the seal state of a Vault server.
type SealStatus struct {
	Sealed      bool   `json:"sealed"`
	Initialized bool   `json:"initialized"`
	T           int    `json:"t"`
	N           int    `json:"n"`
	Progress    int    `json:"progress"`
	Version     string `json:"version"`
	ClusterName string `json:"cluster_name"`
	ClusterID   string `json:"cluster_id"`
}

// GetSealStatus retrieves the current seal status of the Vault server.
func GetSealStatus(client HTTPClient, addr string) (*SealStatus, error) {
	url := fmt.Sprintf("%s/v1/sys/seal-status", addr)

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("seal status request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var status SealStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode seal status: %w", err)
	}

	return &status, nil
}

// IsSealed returns true if the Vault server is currently sealed.
func IsSealed(client HTTPClient, addr string) (bool, error) {
	status, err := GetSealStatus(client, addr)
	if err != nil {
		return false, err
	}
	return status.Sealed, nil
}
