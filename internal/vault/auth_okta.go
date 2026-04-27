package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// OktaLoginResponse holds the Vault token returned after a successful Okta login.
type OktaLoginResponse struct {
	Auth struct {
		ClientToken string `json:"client_token"`
	} `json:"auth"`
}

// OktaLogin authenticates against Vault using the Okta auth method.
// mount defaults to "okta" if empty.
func OktaLogin(c *Client, username, password, mount string) (string, error) {
	if mount == "" {
		mount = "okta"
	}

	payload := map[string]string{
		"password": password,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("okta login: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/v1/auth/%s/login/%s", c.Address, mount, username)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("okta login: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", fmt.Errorf("okta login: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("okta login: unexpected status %d", resp.StatusCode)
	}

	var result OktaLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("okta login: decode response: %w", err)
	}

	if result.Auth.ClientToken == "" {
		return "", fmt.Errorf("okta login: empty token in response")
	}

	return result.Auth.ClientToken, nil
}
