package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// LDAPLoginResponse holds the token returned from an LDAP login.
type LDAPLoginResponse struct {
	Token    string
	Policies []string
	LeaseDur int
}

// LDAPLogin authenticates against Vault using the LDAP auth method.
// mount is typically "ldap" but can be customised.
func LDAPLogin(c *Client, username, password, mount string) (*LDAPLoginResponse, error) {
	if mount == "" {
		mount = "ldap"
	}

	path := fmt.Sprintf("/v1/auth/%s/login/%s", mount, username)
	body, err := jsonBody(map[string]string{"password": password})
	if err != nil {
		return nil, fmt.Errorf("ldap login: encode request: %w", err)
	}

	resp, err := c.RawRequest(http.MethodPost, path, body)
	if err != nil {
		return nil, fmt.Errorf("ldap login: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ldap login: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Auth struct {
			ClientToken string   `json:"client_token"`
			Policies    []string `json:"policies"`
			LeaseDur    int      `json:"lease_duration"`
		} `json:"auth"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("ldap login: decode response: %w", err)
	}
	if result.Auth.ClientToken == "" {
		return nil, fmt.Errorf("ldap login: empty token in response")
	}

	return &LDAPLoginResponse{
		Token:    result.Auth.ClientToken,
		Policies: result.Auth.Policies,
		LeaseDur: result.Auth.LeaseDur,
	}, nil
}
