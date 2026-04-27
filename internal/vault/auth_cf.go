package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// CFLoginRequest holds the parameters for Cloud Foundry authentication.
type CFLoginRequest struct {
	RoleID      string `json:"role"`
	SigningTime string `json:"signing_time"`
	CFInstanceCert string `json:"cf_instance_cert"`
	Signature   string `json:"signature"`
	Mount       string `json:"-"`
}

// CFLoginResponse holds the Vault response for a CF login.
type CFLoginResponse struct {
	Auth struct {
		ClientToken string `json:"client_token"`
	} `json:"auth"`
}

// CFLogin authenticates with Vault using the Cloud Foundry auth method.
func CFLogin(c *Client, req CFLoginRequest) (string, error) {
	mount := req.Mount
	if mount == "" {
		mount = "cf"
	}

	body, err := jsonMarshal(map[string]string{
		"role":             req.RoleID,
		"signing_time":     req.SigningTime,
		"cf_instance_cert": req.CFInstanceCert,
		"signature":        req.Signature,
	})
	if err != nil {
		return "", fmt.Errorf("cf login: marshal request: %w", err)
	}

	resp, err := c.Post(fmt.Sprintf("/v1/auth/%s/login", mount), body)
	if err != nil {
		return "", fmt.Errorf("cf login: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cf login: unexpected status %d", resp.StatusCode)
	}

	var result CFLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("cf login: decode response: %w", err)
	}

	if result.Auth.ClientToken == "" {
		return "", fmt.Errorf("cf login: empty token in response")
	}

	return result.Auth.ClientToken, nil
}
