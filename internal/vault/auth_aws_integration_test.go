package vault_test

import (
	"os"
	"testing"

	"github.com/yourusername/vaultpipe/internal/vault"
)

// TestAWSLogin_Integration runs against a real Vault instance.
// Set VAULT_INTEGRATION=1 along with the required env vars to enable.
func TestAWSLogin_Integration(t *testing.T) {
	if os.Getenv("VAULT_INTEGRATION") != "1" {
		t.Skip("skipping integration test; set VAULT_INTEGRATION=1 to run")
	}

	addr := os.Getenv("VAULT_ADDR")
	token := os.Getenv("VAULT_TOKEN")
	if addr == "" || token == "" {
		t.Fatal("VAULT_ADDR and VAULT_TOKEN must be set for integration tests")
	}

	role := os.Getenv("VAULT_AWS_ROLE")
	iamURL := os.Getenv("VAULT_AWS_IAM_URL")
	iamBody := os.Getenv("VAULT_AWS_IAM_BODY")
	iamHeaders := os.Getenv("VAULT_AWS_IAM_HEADERS")

	if role == "" || iamURL == "" || iamBody == "" {
		t.Skip("skipping: VAULT_AWS_ROLE, VAULT_AWS_IAM_URL, VAULT_AWS_IAM_BODY must be set")
	}

	c, err := vault.NewClient(addr, token)
	if err != nil {
		t.Fatalf("create client: %v", err)
	}

	resultToken, err := vault.AWSLogin(c, role, iamURL, iamBody, iamHeaders, "")
	if err != nil {
		t.Fatalf("AWSLogin: %v", err)
	}
	if resultToken == "" {
		t.Error("expected non-empty token")
	}
	t.Logf("received token: %s", resultToken[:8]+"...")
}
