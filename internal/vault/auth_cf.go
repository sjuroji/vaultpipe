package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// CFLogin authenticates using the Cloud Foundry auth method.
func CFLogin(client *Client, role, cfInstanceCertContents, cfInstanceKeyContents, signingTime, signature string, mount string) (string, error) {
	if mount == "" {
		mount = "cf"
	}

	payload := map[string]string{
		"role":                      role,
		"cf_instance_cert":          cfInstanceCertContents,
		"cf_instance_key":           cfInstanceKeyContents,
		"signing_time":              signingTime,
		"signature":                 signature,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("cf login: marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/v1/auth/%s/login", client.Address, mount)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("cf login: build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.HTTP.Do(req)
	if err != nil {
		return "", fmt.Errorf("cf login: do request: %w", err)
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
		return "", fmt.Errorf("cf login: decode response: %w", err)
	}
	if result.Auth.ClientToken == "" {
		return "", fmt.Errorf("cf login: empty token in response")
	}
	return result.Auth.ClientToken, nil
}
