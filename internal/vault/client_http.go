package vault

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// getJSON performs a GET against the given path and decodes JSON into dest.
func (c *Client) getJSON(ctx context.Context, path string, dest interface{}) error {
	url := fmt.Sprintf("%s/%s", c.addr, path)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Vault-Token", c.token)
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("vault returned status %d for GET %s", resp.StatusCode, path)
	}
	return json.NewDecoder(resp.Body).Decode(dest)
}

// postJSON performs a POST with a JSON body against the given path.
func (c *Client) postJSON(ctx context.Context, path string, body interface{}) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/%s", c.addr, path)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("X-Vault-Token", c.token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("vault returned status %d for POST %s", resp.StatusCode, path)
	}
	return nil
}
