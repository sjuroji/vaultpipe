package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// OIDCLoginResponse holds the token returned after a successful OIDC login.
type OIDCLoginResponse struct {
	Token string
}

// OIDCLogin authenticates against Vault using an OIDC JWT token.
// mount defaults to "oidc" if empty.
func OIDCLogin(c *Client, jwt, role, mount string) (*OIDCLoginResponse, error) {
	if mount == "" {
		mount = "oidc"
	}

	path := fmt.Sprintf("/v1/auth/%s/login", mount)

	body := map[string]string{
		"jwt":  jwt,
		"role": role,
	}

	resp, err := c.PostJSON(path, body)
	if err != nil {
		return nil, fmt.Errorf("oidc login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("oidc login: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Auth struct {
			ClientToken string `json:"client_token"`
		} `json:"auth"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("oidc login decode: %w", err)
	}

	if result.Auth.ClientToken == "" {
		return nil, fmt.Errorf("oidc login: empty token in response")
	}

	return &OIDCLoginResponse{Token: result.Auth.ClientToken}, nil
}
