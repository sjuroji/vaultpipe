package vault

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// SnapshotMeta holds metadata about a Vault raft snapshot.
type SnapshotMeta struct {
	Index     uint64    `json:"index"`
	Term      uint64    `json:"term"`
	Version   int       `json:"version"`
	Timestamp time.Time `json:"timestamp"`
}

// TakeSnapshot requests a raft snapshot from Vault and writes the raw
// snapshot bytes to dst. It returns the number of bytes written.
func TakeSnapshot(c *Client, dst io.Writer) (int64, error) {
	req, err := http.NewRequest(http.MethodGet, c.Address+"/v1/sys/storage/raft/snapshot", nil)
	if err != nil {
		return 0, fmt.Errorf("snapshot: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.Token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return 0, fmt.Errorf("snapshot: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("snapshot: unexpected status %d: %s", resp.StatusCode, body)
	}

	n, err := io.Copy(dst, resp.Body)
	if err != nil {
		return n, fmt.Errorf("snapshot: write: %w", err)
	}
	return n, nil
}

// SnapshotStatus returns metadata about the current raft snapshot state.
func SnapshotStatus(c *Client) (*SnapshotMeta, error) {
	req, err := http.NewRequest(http.MethodGet, c.Address+"/v1/sys/storage/raft/snapshot-auto/status", nil)
	if err != nil {
		return nil, fmt.Errorf("snapshot status: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.Token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("snapshot status: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("snapshot status: unexpected status %d: %s", resp.StatusCode, body)
	}

	var meta SnapshotMeta
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return nil, fmt.Errorf("snapshot status: decode: %w", err)
	}
	return &meta, nil
}
