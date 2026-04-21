package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SecretList holds a list of secret keys at a given path.
type SecretList struct {
	Keys []string
}

// ListKV returns the keys stored under the given KV v2 path.
func ListKV(client *Client, mount, path string) (*SecretList, error) {
	url := fmt.Sprintf("%s/v1/%s/metadata/%s", client.Address, mount, path)

	req, err := http.NewRequest("LIST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("building LIST request: %w", err)
	}
	req.Header.Set("X-Vault-Token", client.Token)

	resp, err := client.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing LIST request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vault returned status %d for LIST %s/%s", resp.StatusCode, mount, path)
	}

	var body struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("decoding LIST response: %w", err)
	}

	return &SecretList{Keys: body.Data.Keys}, nil
}

// DeleteKV deletes the secret at the given KV v2 path (all versions).
func DeleteKV(client *Client, mount, path string) error {
	url := fmt.Sprintf("%s/v1/%s/metadata/%s", client.Address, mount, path)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("building DELETE request: %w", err)
	}
	req.Header.Set("X-Vault-Token", client.Token)

	resp, err := client.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("executing DELETE request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("vault returned status %d for DELETE %s/%s", resp.StatusCode, mount, path)
	}

	return nil
}
