package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// SecretEngine represents a mounted secret engine in Vault.
type SecretEngine struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Accessor    string `json:"accessor"`
	Local       bool   `json:"local"`
	SealWrap    bool   `json:"seal_wrap"`
}

// ListSecretEngines returns all mounted secret engines from /v1/sys/mounts.
func ListSecretEngines(c *Client) (map[string]SecretEngine, error) {
	req, err := http.NewRequest(http.MethodGet, c.Address+"/v1/sys/mounts", nil)
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

	var result map[string]SecretEngine
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}
	return result, nil
}

// GetSecretEngine returns a single secret engine by its mount path.
func GetSecretEngine(c *Client, path string) (SecretEngine, error) {
	engines, err := ListSecretEngines(c)
	if err != nil {
		return SecretEngine{}, err
	}
	// Vault keys include trailing slash
	key := path
	if len(key) == 0 || key[len(key)-1] != '/' {
		key += "/"
	}
	if engine, ok := engines[key]; ok {
		return engine, nil
	}
	return SecretEngine{}, fmt.Errorf("secret engine %q not found", path)
}
