package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Client is a minimal Vault HTTP client.
type Client struct {
	Address string
	Token   string
	http    *http.Client
}

// NewClient creates a Client from environment variables or explicit values.
func NewClient(address, token string) *Client {
	if address == "" {
		address = os.Getenv("VAULT_ADDR")
	}
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}
	return &Client{
		Address: strings.TrimRight(address, "/"),
		Token:   token,
		http:    &http.Client{Timeout: 10 * time.Second},
	}
}

// ReadSecret reads a KV v2 secret at the given path and returns key/value pairs.
func (c *Client) ReadSecret(path string) (map[string]string, error) {
	url := fmt.Sprintf("%s/v1/%s", c.Address, path)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("vault: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.Token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("vault: request failed: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vault: unexpected status %d: %s", resp.StatusCode, body)
	}

	var result struct {
		Data struct {
			Data map[string]interface{} `json:"data"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("vault: parse response: %w", err)
	}

	out := make(map[string]string, len(result.Data.Data))
	for k, v := range result.Data.Data {
		out[k] = fmt.Sprintf("%v", v)
	}
	return out, nil
}
