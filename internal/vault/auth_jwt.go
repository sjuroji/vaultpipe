package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// JWTLoginResponse holds the client token returned after a successful JWT login.
type JWTLoginResponse struct {
	Token string
}

// JWTLogin authenticates against Vault using a JWT/OIDC token.
// mount defaults to "jwt" if empty.
func JWTLogin(client *Client, jwt, role, mount string) (*JWTLoginResponse, error) {
	if mount == "" {
		mount = "jwt"
	}

	body := fmt.Sprintf(`{"jwt":%q,"role":%q}`, jwt, role)
	url := fmt.Sprintf("%s/v1/auth/%s/login", strings.TrimRight(client.Address, "/"), mount)

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("jwt login: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("jwt login: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jwt login: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Auth struct {
			ClientToken string `json:"client_token"`
		} `json:"auth"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("jwt login: decode response: %w", err)
	}
	if result.Auth.ClientToken == "" {
		return nil, fmt.Errorf("jwt login: empty token in response")
	}
	return &JWTLoginResponse{Token: result.Auth.ClientToken}, nil
}
