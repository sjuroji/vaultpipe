package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// UserpassLoginResponse holds the token returned after userpass authentication.
type UserpassLoginResponse struct {
	Token     string
	LeaseDur  int
	Renewable bool
}

// UserpassLogin authenticates with Vault using the userpass auth method and
// returns a client token on success.
func UserpassLogin(c *Client, username, password, mountPath string) (*UserpassLoginResponse, error) {
	if mountPath == "" {
		mountPath = "userpass"
	}

	body := map[string]string{"password": password}
	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("userpass: marshal request: %w", err)
	}

	path := fmt.Sprintf("/v1/auth/%s/login/%s", mountPath, username)
	req, err := http.NewRequest(http.MethodPost, c.Address+path, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("userpass: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("userpass: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("userpass: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Auth struct {
			ClientToken string `json:"client_token"`
			LeaseDuration int  `json:"lease_duration"`
			Renewable   bool   `json:"renewable"`
		} `json:"auth"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("userpass: decode response: %w", err)
	}
	if result.Auth.ClientToken == "" {
		return nil, fmt.Errorf("userpass: empty token in response")
	}

	return &UserpassLoginResponse{
		Token:     result.Auth.ClientToken,
		LeaseDur:  result.Auth.LeaseDuration,
		Renewable: result.Auth.Renewable,
	}, nil
}
