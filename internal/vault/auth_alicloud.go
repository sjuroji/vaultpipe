package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// AliCloudLoginRequest holds credentials for Alibaba Cloud auth.
type AliCloudLoginRequest struct {
	IdentityRequestURL    string `json:"identity_request_url"`
	IdentityRequestHeaders string `json:"identity_request_headers"`
	Role                  string `json:"role"`
}

// AliCloudLoginResponse holds the token returned after a successful login.
type AliCloudLoginResponse struct {
	Token string
}

// AliCloudLogin authenticates against Vault using the AliCloud auth method.
// mount defaults to "alicloud" if empty.
func AliCloudLogin(c *Client, role, identityURL, identityHeaders, mount string) (*AliCloudLoginResponse, error) {
	if mount == "" {
		mount = "alicloud"
	}

	payload := AliCloudLoginRequest{
		IdentityRequestURL:    identityURL,
		IdentityRequestHeaders: identityHeaders,
		Role:                  role,
	}

	body, err := jsonBody(payload)
	if err != nil {
		return nil, fmt.Errorf("alicloud login: encode request: %w", err)
	}

	path := fmt.Sprintf("/v1/auth/%s/login", mount)
	resp, err := c.Post(path, body)
	if err != nil {
		return nil, fmt.Errorf("alicloud login: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("alicloud login: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Auth struct {
			ClientToken string `json:"client_token"`
		} `json:"auth"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("alicloud login: decode response: %w", err)
	}

	if result.Auth.ClientToken == "" {
		return nil, fmt.Errorf("alicloud login: empty token in response")
	}

	return &AliCloudLoginResponse{Token: result.Auth.ClientToken}, nil
}
