package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// RaftConfiguration holds the Raft cluster configuration.
type RaftConfiguration struct {
	Index   int           `json:"index"`
	Servers []RaftServer  `json:"servers"`
}

// RaftServer represents a single node in the Raft cluster.
type RaftServer struct {
	ID       string `json:"node_id"`
	Address  string `json:"address"`
	Leader   bool   `json:"leader"`
	Voter    bool   `json:"voter"`
	Protocol string `json:"protocol_version"`
}

// GetRaftConfiguration retrieves the current Raft cluster configuration.
func GetRaftConfiguration(c *Client) (*RaftConfiguration, error) {
	req, err := http.NewRequest(http.MethodGet, c.Address+"/v1/sys/storage/raft/configuration", nil)
	if err != nil {
		return nil, fmt.Errorf("raft config request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.Token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("raft config: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("raft config: unexpected status %d", resp.StatusCode)
	}

	var wrapper struct {
		Data RaftConfiguration `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wrapper); err != nil {
		return nil, fmt.Errorf("raft config decode: %w", err)
	}
	return &wrapper.Data, nil
}

// RemoveRaftPeer removes a peer from the Raft cluster by node ID.
func RemoveRaftPeer(c *Client, nodeID string) error {
	body, err := jsonBody(map[string]string{"server_id": nodeID})
	if err != nil {
		return fmt.Errorf("raft remove peer body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.Address+"/v1/sys/storage/raft/remove-peer", body)
	if err != nil {
		return fmt.Errorf("raft remove peer request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return fmt.Errorf("raft remove peer: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("raft remove peer: unexpected status %d", resp.StatusCode)
	}
	return nil
}
