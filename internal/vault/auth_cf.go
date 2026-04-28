package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// CFLoginRequest holds credentials for Cloud Foundry authentication.
type CFLoginRequest struct {
	RoleID      string `json:"role"`
	SigningTime string `json:"signing_time"`
	CFInstanceCert string `json:"cf_instance_cert"`
	Signature   string `json:"signature"`
}

// CFLogin authenticates using the Cloud Foundry auth method.
// mount defaults to "cf" if empty.
func CFLogin(c *Client, mount, role, signingTime, cert, signature string) (string, error) {
	if mount == "" {
		mount = "cf"
	}

	payload := CFLoginRequest{
		RoleID:         role,
		SigningTime:    signingTime,
		CFInstanceCert: cert,
		Signature:      signature,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("cf login: marshal: %w", err)
	}

	url := fmt.Sprintf("%s/v1/auth/%s/login", c.Address, mount)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("cf login: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", fmt.Errorf("cf login: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cf login: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Auth struct {
			ClientToken string `json:"client_token"`
		} `json:"auth"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("cf login: decode: %w", err)
	}

	if result.Auth.ClientToken == "" {
		return "", fmt.Errorf("cf login: empty token in response")
	}

	return result.Auth.ClientToken, nil
}
