package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// OCILoginRequest holds the credentials for OCI IAM authentication.
type OCILoginRequest struct {
	RequestHeaders map[string]string `json:"request_headers"`
}

// OCILoginResponse holds the parsed token from an OCI login response.
type OCILoginResponse struct {
	Token string
}

// OCILogin authenticates with Vault using OCI IAM and returns a client token.
// mount defaults to "oci" if empty.
func OCILogin(c *Client, role string, requestHeaders map[string]string, mount string) (*OCILoginResponse, error) {
	if mount == "" {
		mount = "oci"
	}

	payload := OCILoginRequest{
		RequestHeaders: requestHeaders,
	}

	body, err := jsonMarshal(payload)
	if err != nil {
		return nil, fmt.Errorf("oci login: marshal payload: %w", err)
	}

	path := fmt.Sprintf("/v1/auth/%s/login/%s", mount, role)
	resp, err := c.RawPost(path, body)
	if err != nil {
		return nil, fmt.Errorf("oci login: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("oci login: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Auth struct {
			ClientToken string `json:"client_token"`
		} `json:"auth"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("oci login: decode response: %w", err)
	}

	if result.Auth.ClientToken == "" {
		return nil, fmt.Errorf("oci login: empty token in response")
	}

	return &OCILoginResponse{Token: result.Auth.ClientToken}, nil
}
