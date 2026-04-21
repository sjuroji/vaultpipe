package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Plugin represents a registered Vault plugin.
type Plugin struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Version string `json:"version"`
	Builtin bool   `json:"builtin"`
}

// PluginListResponse is the response from the catalog list endpoint.
type PluginListResponse struct {
	Data struct {
		Auth     []string `json:"auth"`
		Database []string `json:"database"`
		Secret   []string `json:"secret"`
	} `json:"data"`
}

// ListPlugins returns all registered plugins from the Vault plugin catalog.
func ListPlugins(c *Client) (*PluginListResponse, error) {
	req, err := http.NewRequest(http.MethodGet, c.Address+"/v1/sys/plugins/catalog", nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.Token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var result PluginListResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return &result, nil
}

// GetPlugin retrieves details for a specific plugin by type and name.
func GetPlugin(c *Client, pluginType, name string) (*Plugin, error) {
	url := fmt.Sprintf("%s/v1/sys/plugins/catalog/%s/%s", c.Address, pluginType, name)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.Token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("plugin %q of type %q not found", name, pluginType)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var wrapper struct {
		Data Plugin `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return &wrapper.Data, nil
}
