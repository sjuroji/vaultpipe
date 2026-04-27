package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// AzureLoginRequest holds the credentials for Azure auth.
type AzureLoginRequest struct {
	Role           string `json:"role"`
	JWT            string `json:"jwt"`
	SubscriptionID string `json:"subscription_id,omitempty"`
	ResourceGroup  string `json:"resource_group_name,omitempty"`
	VMName         string `json:"vm_name,omitempty"`
	VMSSName       string `json:"vmss_name,omitempty"`
}

// AzureLogin authenticates using the Azure auth method and returns a Vault token.
func AzureLogin(c *Client, role, jwt, mount string) (string, error) {
	if mount == "" {
		mount = "azure"
	}

	payload := AzureLoginRequest{
		Role: role,
		JWT:  jwt,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("azure login: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/v1/auth/%s/login", c.Address, mount)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("azure login: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", fmt.Errorf("azure login: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("azure login: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Auth struct {
			ClientToken string `json:"client_token"`
		} `json:"auth"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("azure login: decode response: %w", err)
	}

	if result.Auth.ClientToken == "" {
		return "", fmt.Errorf("azure login: empty token in response")
	}

	return result.Auth.ClientToken, nil
}
