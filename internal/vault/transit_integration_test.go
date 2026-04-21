package vault

import (
	"encoding/base64"
	"os"
	"testing"
)

// TestTransit_Integration performs a real encrypt→decrypt round-trip.
// Requires: VAULT_ADDR, VAULT_TOKEN, VAULT_TRANSIT_KEY env vars and a
// running Vault instance with the transit engine enabled.
func TestTransit_Integration(t *testing.T) {
	addr := os.Getenv("VAULT_ADDR")
	token := os.Getenv("VAULT_TOKEN")
	key := os.Getenv("VAULT_TRANSIT_KEY")
	if addr == "" || token == "" || key == "" {
		t.Skip("skipping integration test: VAULT_ADDR, VAULT_TOKEN, VAULT_TRANSIT_KEY not set")
	}

	c, err := NewClient(addr, token)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	original := "super-secret-value"
	b64 := base64.StdEncoding.EncodeToString([]byte(original))

	encRes, err := EncryptTransit(c, key, b64)
	if err != nil {
		t.Fatalf("EncryptTransit: %v", err)
	}
	if encRes.Ciphertext == "" {
		t.Fatal("expected non-empty ciphertext")
	}

	decRes, err := DecryptTransit(c, key, encRes.Ciphertext)
	if err != nil {
		t.Fatalf("DecryptTransit: %v", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(decRes.Plaintext)
	if err != nil {
		t.Fatalf("base64 decode: %v", err)
	}
	if string(decoded) != original {
		t.Errorf("round-trip mismatch: got %q, want %q", string(decoded), original)
	}
}
