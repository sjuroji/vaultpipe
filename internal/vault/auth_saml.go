package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// SAMLLoginRequest holds the parameters for SAML authentication.
type SAMLLoginRequest struct {
	SAMLResponse string `json:"saml_response"`
	Mount        string `json:"-"`
}

// SAMLLoginResponse holds the result of a successful SAML login.
type SAMLLoginResponse struct {
	ClientToken string
	LeaseDuration int
	Renewable     bool
}

// SAMLLogin authenticates against Vault using a SAML assertion.
// mount defaults to "saml" if empty.
func SAMLLogin(c *Client, req SAMLLoginRequest) (*SAMLLoginResponse, error) {
	mount := req.Mount
	if mount == "" {
		mount = "saml"
	}

	body := map[string]string{
		"saml_response": req.SAMLResponse,
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("saml: marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/v1/auth/%s/login", c.Address, mount)
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("saml: build request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if c.Token != "" {
		httpReq.Header.Set("X-Vault-Token", c.Token)
	}

	resp, err := c.HTTP.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("saml: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("saml: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Auth struct {
			ClientToken   string `json:"client_token"`
			LeaseDuration int    `json:"lease_duration"`
			Renewable     bool   `json:"renewable"`
		} `json:"auth"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("saml: decode response: %w", err)
	}
	if result.Auth.ClientToken == "" {
		return nil, fmt.Errorf("saml: empty client token in response")
	}

	return &SAMLLoginResponse{
		ClientToken:   result.Auth.ClientToken,
		LeaseDuration: result.Auth.LeaseDuration,
		Renewable:     result.Auth.Renewable,
	}, nil
}
