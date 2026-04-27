package vault

import (
	"fmt"
	"net/http"
)

// TLSLoginResponse holds the token returned from a TLS certificate login.
type TLSLoginResponse struct {
	Token string
	Mount string
}

// TLSLogin authenticates to Vault using a TLS client certificate via the
// cert auth method. The mount parameter defaults to "cert" if empty.
// Unlike CertLogin (which sends explicit cert/key paths), TLSLogin relies on
// the TLS client certificate already configured on the underlying HTTP client.
func TLSLogin(c *Client, mount string) (*TLSLoginResponse, error) {
	if mount == "" {
		mount = "cert"
	}

	path := fmt.Sprintf("/v1/auth/%s/login", mount)

	resp, err := c.RawRequest(http.MethodPost, path, nil)
	if err != nil {
		return nil, fmt.Errorf("tls login request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tls login: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Auth struct {
			ClientToken string `json:"client_token"`
		} `json:"auth"`
	}
	if err := decodeJSON(resp.Body, &result); err != nil {
		return nil, fmt.Errorf("tls login: decode response: %w", err)
	}

	if result.Auth.ClientToken == "" {
		return nil, fmt.Errorf("tls login: empty token in response")
	}

	return &TLSLoginResponse{
		Token: result.Auth.ClientToken,
		Mount: mount,
	}, nil
}
