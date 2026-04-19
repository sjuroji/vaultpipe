package vault

import (
	"testing"
	"time"
)

func TestParseLease_OK(t *testing.T) {
	data := map[string]interface{}{
		"lease_id":       "secret/data/foo/bar",
		"lease_duration": float64(3600),
		"renewable":      true,
	}

	lease, err := ParseLease(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lease.LeaseID != "secret/data/foo/bar" {
		t.Errorf("expected lease_id %q, got %q", "secret/data/foo/bar", lease.LeaseID)
	}
	if lease.LeaseDuration != 3600*time.Second {
		t.Errorf("expected duration 3600s, got %v", lease.LeaseDuration)
	}
	if !lease.Renewable {
		t.Error("expected renewable to be true")
	}
}

func TestParseLease_MissingDuration(t *testing.T) {
	data := map[string]interface{}{
		"lease_id":  "abc",
		"renewable": false,
	}
	_, err := ParseLease(data)
	if err == nil {
		t.Fatal("expected error for missing lease_duration")
	}
}

func TestParseLease_BadDurationType(t *testing.T) {
	data := map[string]interface{}{
		"lease_id":       "abc",
		"lease_duration": "not-a-number",
		"renewable":      false,
	}
	_, err := ParseLease(data)
	if err == nil {
		t.Fatal("expected error for bad lease_duration type")
	}
}

func TestIsExpiringSoon(t *testing.T) {
	lease := &LeaseInfo{LeaseDuration: 30 * time.Second}
	if !lease.IsExpiringSoon(1 * time.Minute) {
		t.Error("expected lease to be expiring soon")
	}
	if lease.IsExpiringSoon(10 * time.Second) {
		t.Error("expected lease NOT to be expiring soon")
	}
}
