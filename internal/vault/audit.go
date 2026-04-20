package vault

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AuditEvent represents a single secret access event.
type AuditEvent struct {
	Path      string    `json:"path"`
	Operation string    `json:"operation"`
	Timestamp time.Time `json:"timestamp"`
	Success   bool      `json:"success"`
}

// AuditLog holds a slice of AuditEvents.
type AuditLog struct {
	Events []AuditEvent
}

// Record appends a new event to the audit log.
func (a *AuditLog) Record(path, operation string, success bool) {
	a.Events = append(a.Events, AuditEvent{
		Path:      path,
		Operation: operation,
		Timestamp: time.Now().UTC(),
		Success:   success,
	})
}

// AuditResponse is the Vault audit backends list response.
type AuditResponse struct {
	Data map[string]AuditBackend `json:"data"`
}

// AuditBackend describes a single audit backend.
type AuditBackend struct {
	Type        string `json:"type"`
	Description string `json:"description"`
}

// ListAuditBackends retrieves the configured audit backends from Vault.
func ListAuditBackends(c *Client) (map[string]AuditBackend, error) {
	req, err := http.NewRequest(http.MethodGet, c.Address+"/v1/sys/audit", nil)
	if err != nil {
		return nil, fmt.Errorf("audit: build request: %w", err)
	}
	req.Header.Set("X-Vault-Token", c.Token)

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("audit: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("audit: unexpected status %d", resp.StatusCode)
	}

	var ar AuditResponse
	if err := json.NewDecoder(resp.Body).Decode(&ar); err != nil {
		return nil, fmt.Errorf("audit: decode response: %w", err)
	}

	return ar.Data, nil
}
