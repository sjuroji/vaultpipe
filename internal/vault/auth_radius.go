package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RADIUSLoginResponse represents the response from a RADIUS authentication request.
type RADIUSLoginResponse struct {
	Auth struct {
		ClientToken string `json:"client_token"`
	} `json:"auth"`
}

// RADIUSLogin authenticates against Vault using the RADIUS auth method.
// mount specifies the mount path (defaults to "radius" if empty).
// username and password are the RADIUS credentials.
// Returns the client token on success.
func RADIUSLogin(c *Client, username, password, mount string) (string, error) {
	if mount == "" {
		mount = "radius"
	}

	if username == "" {
		return "", fmt.Errorf("radius login: username is required")
	}
	if password == "" {
		return "", fmt.Errorf("radius login: password is required")
	}

	path := fmt.Sprintf("/v1/auth/%s/login/%s", mount, username)

	body := map[string]string{
		"password": password,
	}

	resp, err := c.PostJSON(path, body)
	if err != nil {
		return "", fmt.Errorf("radius login: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("radius login: unexpected status %d", resp.StatusCode)
	}

	var result RADIUSLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("radius login: failed to decode response: %w", err)
	}

	if result.Auth.ClientToken == "" {
		return "", fmt.Errorf("radius login: empty token in response")
	}

	return result.Auth.ClientToken, nil
}
