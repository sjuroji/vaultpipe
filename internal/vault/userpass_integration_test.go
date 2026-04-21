//go:build integration
// +build integration

package vault_test

import (
	"os"
	"testing"

	"github.com/yourorg/vaultpipe/internal/vault"
)

// requires a live Vault instance with userpass
// auth enabled. Set VAULT_ADDR, VAULT_TEST_USER, and VAULT_TEST_PASS.
func TestUserpassLogin_Integration(t *testing.T) {
	addr := os.Getenv("VAULT_ADDR")
	user := os.Getenv("VAULT_TEST_USER")
	pass := os.Getenv("VAULT_TEST_PASS")

	if addr == "" || user == "" || pass == "" {
		t.Skip("VAULT_ADDR, VAULT_TEST_USER, VAULT_TEST_PASS not set")
	}

	client, err := vault.NewClient(addr, "")
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}

	resp, err := vault.UserpassLogin(client, user, pass, "")
	if err != nil {
		t.Fatalf("UserpassLogin: %v", err)
	}
	if resp.Token == "" {
		t.Error("expected non-empty token")
	}
	t.Logf("token: %s lease: %ds renewable: %v", resp.Token, resp.LeaseDur, resp.Renewable)
}
