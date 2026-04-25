package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// AppRoleLoginRequest holds the credentials for AppRole authentication.
type AppRoleLoginRequest struct {
	RoleID   string `json:"role_id"`
	SecretID string `json:"secret_id"`
}

// AppRoleLoginResponse holds the Vault response for an AppRole login.
type AppRoleLoginResponse struct {
	Auth struct {
		ClientToken string `json:"client_token"`
		LeaseDuration int    `json:"lease_duration"`
		Renewable     bool   `json:"renewable"`
	} `json:"auth"`
}

// AppRoleLogin authenticates with Vault using the AppRole method and returns
// the client token. mount defaults to "approle" if empty.
func AppRoleLogin(c *Client, roleID, secretID, mount string) (string, error) {
	if mount == "" {
		mount = "approle"
	}

	payload, err := json.Marshal(AppRoleLoginRequest{
		RoleID:   roleID,
		SecretID: secretID,
	})
	if err != nil {
		return "", fmt.Errorf("approle: marshal request: %w", err)
	}

	path := fmt.Sprintf("/v1/auth/%s/login", strings.Trim(mount, "/"))
	resp, err := c.Post(path, payload)
	if err != nil {
		return "", fmt.Errorf("approle: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("approle: unexpected status %d", resp.StatusCode)
	}

	var result AppRoleLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("approle: decode response: %w", err)
	}

	if result.Auth.ClientToken == "" {
		return "", fmt.Errorf("approle: empty token in response")
	}

	return result.Auth.ClientToken, nil
}
