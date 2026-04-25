package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// GCPLoginRequest holds the parameters for GCP IAM authentication.
type GCPLoginRequest struct {
	Role      string `json:"role"`
	JWT       string `json:"jwt"`
	Mount     string `json:"-"`
}

// GCPLoginResponse holds the client token returned after a successful GCP login.
type GCPLoginResponse struct {
	ClientToken string
}

// GCPLogin authenticates against Vault using a GCP IAM JWT and returns a client token.
func GCPLogin(c *Client, req GCPLoginRequest) (*GCPLoginResponse, error) {
	mount := req.Mount
	if mount == "" {
		mount = "gcp"
	}

	body := map[string]string{
		"role": req.Role,
		"jwt":  req.JWT,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("gcp login: marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/auth/%s/login", c.Address, mount)
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("gcp login: create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("gcp login: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gcp login: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Auth struct {
			ClientToken string `json:"client_token"`
		} `json:"auth"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("gcp login: decode response: %w", err)
	}

	if result.Auth.ClientToken == "" {
		return nil, fmt.Errorf("gcp login: empty token in response")
	}

	return &GCPLoginResponse{ClientToken: result.Auth.ClientToken}, nil
}
