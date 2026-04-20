package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// NamespaceInfo holds metadata about a Vault namespace.
type NamespaceInfo struct {
	Path   string            `json:"path"`
	ID     string            `json:"id"`
	Meta   map[string]string `json:"custom_metadata"`
}

// ListNamespaces returns the child namespaces under the given parent path.
// Pass an empty string to list top-level namespaces.
func ListNamespaces(c *Client, parent string) ([]string, error) {
	base := "v1/sys/namespaces"
	if parent != "" {
		base = fmt.Sprintf("v1/sys/namespaces/%s", parent)
	}
	url := fmt.Sprintf("%s/%s", c.Address, base)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("namespace list: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.Token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("namespace list: do request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("namespace list: unexpected status %d: %s", resp.StatusCode, body)
	}

	var payload struct {
		Data struct {
			Keys []string `json:"keys"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("namespace list: decode: %w", err)
	}
	return payload.Data.Keys, nil
}

// GetNamespace returns metadata for the namespace at the given path.
func GetNamespace(c *Client, path string) (*NamespaceInfo, error) {
	url := fmt.Sprintf("%s/v1/sys/namespaces/%s", c.Address, path)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("namespace get: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.Token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("namespace get: do request: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("namespace get: unexpected status %d: %s", resp.StatusCode, body)
	}

	var info NamespaceInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("namespace get: decode: %w", err)
	}
	return &info, nil
}
