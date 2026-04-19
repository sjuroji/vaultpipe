package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AppRoleLoginResponse holds the auth block returned by Vault on AppRole login.
type AppRoleLoginResponse struct {
	Auth struct {
		ClientToken   string `json:"client_token"`
		LeaseDuration int    `json:"lease_duration"`
		Renewable     bool   `json:"renewable"`
	} `json:"auth"`
}

// AppRoleLogin authenticates with Vault using a RoleID and SecretID,
// returning the client token on success.
func AppRoleLogin(c *Client, roleID, secretID string) (string, error) {
	payload := map[string]string{
		"role_id":   roleID,
		"secret_id": secretID,
	}

	body, err := jsonBody(payload)
	if err != nil {
		return "", fmt.Errorf("auth: encode payload: %w", err)
	}

	url := c.BaseURL + "/v1/auth/approle/login"
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return "", fmt.Errorf("auth: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", fmt.Errorf("auth: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("auth: unexpected status %d", resp.StatusCode)
	}

	var result AppRoleLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("auth: decode response: %w", err)
	}

	if result.Auth.ClientToken == "" {
		return "", fmt.Errorf("auth: empty client token in response")
	}

	return result.Auth.ClientToken, nil
}
