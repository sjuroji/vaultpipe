package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// MountInfo holds metadata about a single secrets engine mount.
type MountInfo struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Accessor    string `json:"accessor"`
	Local       bool   `json:"local"`
	SealWrap    bool   `json:"seal_wrap"`
}

// MountsResponse maps mount paths to their MountInfo.
type MountsResponse map[string]MountInfo

// ListMounts returns all secrets engine mounts from Vault.
func ListMounts(c *Client) (MountsResponse, error) {
	req, err := http.NewRequest(http.MethodGet, c.Address+"/v1/sys/mounts", nil)
	if err != nil {
		return nil, fmt.Errorf("vault: build mounts request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.Token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("vault: mounts request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("vault: mounts returned status %d", resp.StatusCode)
	}

	var result MountsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("vault: decode mounts response: %w", err)
	}
	return result, nil
}

// GetMount returns info for a single mount path (e.g. "secret/").
func GetMount(c *Client, path string) (MountInfo, error) {
	mounts, err := ListMounts(c)
	if err != nil {
		return MountInfo{}, err
	}
	info, ok := mounts[path]
	if !ok {
		return MountInfo{}, fmt.Errorf("vault: mount %q not found", path)
	}
	return info, nil
}
