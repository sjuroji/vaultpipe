package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WrappedSecret holds the wrapping token and its metadata.
type WrappedSecret struct {
	Token          string
	CreationTime   time.Time
	TTL            time.Duration
	CreationPath   string
}

// WrapSecret requests Vault to wrap the secret at the given KV path
// using the specified TTL (e.g. "30s", "5m").
func WrapSecret(c *Client, path, wrapTTL string) (*WrappedSecret, error) {
	req, err := http.NewRequest(http.MethodGet, c.Address+"/v1/"+path, nil)
	if err != nil {
		return nil, fmt.Errorf("wrapping: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.Token)
	req.Header.Set("X-Vault-Wrap-TTL", wrapTTL)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("wrapping: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wrapping: unexpected status %d", resp.StatusCode)
	}

	var body struct {
		WrapInfo struct {
			Token        string `json:"token"`
			TTL          int    `json:"ttl"`
			CreationTime string `json:"creation_time"`
			CreationPath string `json:"creation_path"`
		} `json:"wrap_info"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("wrapping: decode response: %w", err)
	}
	if body.WrapInfo.Token == "" {
		return nil, fmt.Errorf("wrapping: empty wrapping token in response")
	}

	ct, err := time.Parse(time.RFC3339Nano, body.WrapInfo.CreationTime)
	if err != nil {
		ct = time.Time{}
	}

	return &WrappedSecret{
		Token:        body.WrapInfo.Token,
		CreationTime: ct,
		TTL:          time.Duration(body.WrapInfo.TTL) * time.Second,
		CreationPath: body.WrapInfo.CreationPath,
	}, nil
}

// UnwrapSecret exchanges a wrapping token for the underlying secret data.
func UnwrapSecret(c *Client, wrappingToken string) (map[string]string, error) {
	req, err := http.NewRequest(http.MethodPost, c.Address+"/v1/sys/wrapping/unwrap", nil)
	if err != nil {
		return nil, fmt.Errorf("unwrap: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", wrappingToken)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unwrap: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unwrap: unexpected status %d", resp.StatusCode)
	}

	var body struct {
		Data map[string]interface{} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, fmt.Errorf("unwrap: decode response: %w", err)
	}

	result := make(map[string]string, len(body.Data))
	for k, v := range body.Data {
		if s, ok := v.(string); ok {
			result[k] = s
		}
	}
	return result, nil
}
