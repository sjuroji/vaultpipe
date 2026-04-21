package vault

import (
	"os"
	"testing"
)

// TestRekey_Integration runs against a real Vault dev server.
// Set VAULT_INTEGRATION=1, VAULT_ADDR, and VAULT_TOKEN to enable.
func TestRekey_Integration(t *testing.T) {
	if os.Getenv("VAULT_INTEGRATION") != "1" {
		t.Skip("skipping integration test; set VAULT_INTEGRATION=1 to run")
	}

	addr := os.Getenv("VAULT_ADDR")
	token := os.Getenv("VAULT_TOKEN")
	if addr == "" || token == "" {
		t.Fatal("VAULT_ADDR and VAULT_TOKEN must be set")
	}

	c, err := NewClient(addr, token)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	// Start a rekey operation.
	status, err := InitRekey(c, 1, 1)
	if err != nil {
		t.Fatalf("InitRekey: %v", err)
	}
	if !status.Started {
		t.Error("expected rekey to be started")
	}
	if status.Nonce == "" {
		t.Error("expected non-empty nonce")
	}

	// Verify status endpoint reflects the started operation.
	got, err := GetRekeyStatus(c)
	if err != nil {
		t.Fatalf("GetRekeyStatus: %v", err)
	}
	if !got.Started {
		t.Error("expected status to show started=true")
	}

	// Cancel to clean up.
	if err := CancelRekey(c); err != nil {
		t.Fatalf("CancelRekey: %v", err)
	}

	// Confirm cancelled.
	after, err := GetRekeyStatus(c)
	if err != nil {
		t.Fatalf("GetRekeyStatus after cancel: %v", err)
	}
	if after.Started {
		t.Error("expected rekey to be cancelled (started=false)")
	}
}
