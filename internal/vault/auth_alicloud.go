package vault

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// AliCloudLogin authenticates using the AliCloud auth method.
func AliCloudLogin(c *Client, role, mount string) (string, error) {
	if mount == "" {
		mount = "alicloud"
	}

	body := map[string]string{
		"role": role,
	}
	b, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("alicloud login: marshal: %w", err)
	}

	url := fmt.Sprintf("%s/v1/auth/%s/login", c.Address, mount)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return "", fmt.Errorf("alicloud login: new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", fmt.Errorf("alicloud login: do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("alicloud login: unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Auth struct {
			ClientToken string `json:"client_token"`
		} `json:"auth"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("alicloud login: decode: %w", err)
	}
	if result.Auth.ClientToken == "" {
		return "", fmt.Errorf("alicloud login: empty token in response")
	}
	return result.Auth.ClientToken, nil
}
