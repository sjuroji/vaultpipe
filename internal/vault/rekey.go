package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// RekeyStatus holds the current rekey operation status.
type RekeyStatus struct {
	Started          bool     `json:"started"`
	Nonce            string   `json:"nonce"`
	T                int      `json:"t"`
	N                int      `json:"n"`
	Progress         int      `json:"progress"`
	Required         int      `json:"required"`
	PGPFingerprints  []string `json:"pgp_fingerprints"`
	Backup           bool     `json:"backup"`
	VerificationRequired bool  `json:"verification_required"`
}

// InitRekey starts a new rekey operation.
func InitRekey(c *Client, secretShares, secretThreshold int) (*RekeyStatus, error) {
	body := map[string]interface{}{
		"secret_shares":    secretShares,
		"secret_threshold": secretThreshold,
	}
	b, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, c.Address+"/v1/sys/rekey/init", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", c.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rekey init failed: status %d", resp.StatusCode)
	}
	var status RekeyStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, err
	}
	return &status, nil
}

// GetRekeyStatus returns the current rekey operation status.
func GetRekeyStatus(c *Client) (*RekeyStatus, error) {
	req, err := http.NewRequest(http.MethodGet, c.Address+"/v1/sys/rekey/init", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Vault-Token", c.Token)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rekey status failed: status %d", resp.StatusCode)
	}
	var status RekeyStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, err
	}
	return &status, nil
}

// CancelRekey cancels any in-progress rekey operation.
func CancelRekey(c *Client) error {
	req, err := http.NewRequest(http.MethodDelete, c.Address+"/v1/sys/rekey/init", nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Vault-Token", c.Token)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("rekey cancel failed: status %d", resp.StatusCode)
	}
	return nil
}
