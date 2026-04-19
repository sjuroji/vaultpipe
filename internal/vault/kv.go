package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// SecretData holds the key-value pairs from a KV v2 secret.
type SecretData map[string]string

// ReadKV reads a KV v2 secret at the given mount and path.
func (c *Client) ReadKV(mount, secretPath string) (SecretData, error) {
	url := fmt.Sprintf("%s/v1/%s/data/%s", c.address, mount, secretPath)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
 resp.StatusCode != http.Statustreturn nil, fmt.Err %s", resp.StatusCode result struct {
		Data struct {
			Data map[string]interface{} `json:"data"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	secrets := make(SecretData, len(result.Data.Data))
	for k, v := range result.Data.Data {
		switch val := v.(type) {
		case string:
			secrets[k] = val
		default:
			secrets[k] = fmt.Sprintf("%v", val)
		}
	}
	return secrets, nil
}
