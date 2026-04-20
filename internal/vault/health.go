package vault

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// HealthStatus represents the response from Vault's /v1/sys/health endpoint.
type HealthStatus struct {
	Initialized                bool   `json:"initialized"`
	Sealed                     bool   `json:"sealed"`
	Standby                    bool   `json:"standby"`
	PerformanceStandby         bool   `json:"performance_standby"`
	ReplicationPerformanceMode string `json:"replication_performance_mode"`
	ReplicationDRMode          string `json:"replication_dr_mode"`
	ServerTimeUTC              int64  `json:"server_time_utc"`
	Version                    string `json:"version"`
	ClusterName                string `json:"cluster_name"`
	ClusterID                  string `json:"cluster_id"`
}

// CheckHealth queries the Vault health endpoint and returns the status.
// Vault returns non-200 codes for certain states (sealed, standby), so
// we treat 200, 429, 472, 473, and 501 as valid readable responses.
func CheckHealth(ctx context.Context, c HTTPClient, addr string) (*HealthStatus, error) {
	url := fmt.Sprintf("%s/v1/sys/health?standbyok=true&perfstandbyok=true", addr)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("vault health: build request: %w", err)
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("vault health: request failed: %w", err)
	}
	defer resp.Body.Close()

	// Vault uses non-200 status codes to signal state, but body is still valid JSON.
	switch resp.StatusCode {
	case http.StatusOK, 429, 472, 473, 501:
		// readable states
	default:
		return nil, fmt.Errorf("vault health: unexpected status %d", resp.StatusCode)
	}

	var status HealthStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("vault health: decode response: %w", err)
	}

	return &status, nil
}

// IsHealthy returns true when Vault is initialized, unsealed, and active.
func IsHealthy(s *HealthStatus) bool {
	return s.Initialized && !s.Sealed && !s.Standby
}
