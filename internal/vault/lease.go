package vault

import (
	"fmt"
	"time"
)

// LeaseInfo holds metadata about a secret lease.
type LeaseInfo struct {
	LeaseID       string
	LeaseDuration time.Duration
	Renewable     bool
}

// ParseLease extracts lease information from a raw Vault API response map.
func ParseLease(data map[string]interface{}) (*LeaseInfo, error) {
	leaseID, _ := data["lease_id"].(string)
	renewable, _ := data["renewable"].(bool)

	durationRaw, ok := data["lease_duration"]
	if !ok {
		return nil, fmt.Errorf("missing lease_duration in response")
	}

	var seconds float64
	switch v := durationRaw.(type) {
	case float64:
		seconds = v
	case int:
		seconds = float64(v)
	default:
		return nil, fmt.Errorf("unexpected type for lease_duration: %T", durationRaw)
	}

	return &LeaseInfo{
		LeaseID:       leaseID,
		LeaseDuration: time.Duration(seconds) * time.Second,
		Renewable:     renewable,
	}, nil
}

// IsExpiringSoon returns true if the lease expires within the given threshold.
func (l *LeaseInfo) IsExpiringSoon(threshold time.Duration) bool {
	return l.LeaseDuration <= threshold
}

// RemainingRatio returns the fraction of the lease duration remaining relative
// to the given total duration. This is useful for determining when to renew,
// for example triggering renewal when the ratio drops below 0.5 (half-life).
// Returns 0 if total is zero to avoid division by zero.
func (l *LeaseInfo) RemainingRatio(total time.Duration) float64 {
	if total <= 0 {
		return 0
	}
	return float64(l.LeaseDuration) / float64(total)
}
