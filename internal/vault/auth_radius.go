package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// RADIUSLoginRequest holds credentials for RADIUS authentication.
type RADIUSLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// RADIUSLoginResponse holds the Vault response from a RADIUS login.
type RADIUSLoginResponse struct {
	Auth struct {
		ClientToken string `json:"client_token"`
	} `json:"auth"`
}

// RADIUSLogin authenticates against Vault using the RADIUS auth method.
// mount defaults to "radius" if empty.
func RADIUSLogin(c *Client, username, password, mount string) (string, error) {
	if mount == "" {
		mount = "radius"
	}

	payload := RADIUSLoginRequest{
		Username: username,
		Password: password,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("radius login: marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/auth/%s/login/%s", c.Address, mount, username)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("radius login: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		req.Header.Set("X-Vault-Token", c.Token)
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", fmt.Errorf("radius login: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("radius login: unexpected status %d", resp.StatusCode)
	}

	var result RADIUSLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("radius login: decode response: %w", err)
	}

	if result.Auth.ClientToken == "" {
		return "", fmt.Errorf("radius login: empty client token in response")
	}

	return result.Auth.ClientToken, nil
}
