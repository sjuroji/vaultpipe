package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// CertLoginResponse holds the token returned from a TLS certificate login.
type CertLoginResponse struct {
	Token     string
	LeaseDur  int
	Renewable bool
}

// CertLogin authenticates to Vault using a TLS certificate at the given mount
// (defaults to "cert" if empty). The client must already be configured with
// the appropriate client certificate via its underlying http.Client.
func CertLogin(c *Client, mount, certRole string) (*CertLoginResponse, error) {
	if mount == "" {
		mount = "cert"
	}

	path := fmt.Sprintf("/v1/auth/%s/login", mount)

	body := map[string]string{}
	if certRole != "" {
		body["name"] = certRole
	}

	resp, err := c.PostJSON(path, body)
	if err != nil {
		return nil, fmt.Errorf("cert login request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cert login returned status %d", resp.StatusCode)
	}

	var result struct {
		Auth struct {
			ClientToken string `json:"client_token"`
			LeaseDuration int  `json:"lease_duration"`
			Renewable   bool   `json:"renewable"`
		} `json:"auth"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("cert login decode failed: %w", err)
	}
	if result.Auth.ClientToken == "" {
		return nil, fmt.Errorf("cert login returned empty token")
	}
	return &CertLoginResponse{
		Token:     result.Auth.ClientToken,
		LeaseDur:  result.Auth.LeaseDuration,
		Renewable: result.Auth.Renewable,
	}, nil
}
