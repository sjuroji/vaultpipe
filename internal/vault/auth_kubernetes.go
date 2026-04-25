package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// KubernetesLogin authenticates using the Kubernetes auth method and returns a Vault token.
// mount defaults to "kubernetes" if empty.
func KubernetesLogin(c *Client, role, jwt, mount string) (string, error) {
	if mount == "" {
		mount = "kubernetes"
	}

	body, err := json.Marshal(map[string]string{
		"role": role,
		"jwt":  jwt,
	})
	if err != nil {
		return "", fmt.Errorf("kubernetes login: marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/auth/%s/login", c.Address, mount)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("kubernetes login: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", fmt.Errorf("kubernetes login: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("kubernetes login: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Auth struct {
			ClientToken string `json:"client_token"`
		} `json:"auth"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("kubernetes login: decode response: %w", err)
	}
	if result.Auth.ClientToken == "" {
		return "", fmt.Errorf("kubernetes login: empty token in response")
	}
	return result.Auth.ClientToken, nil
}
