package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultpipe/internal/vault"
)

func kvV2Response(t *testing.T, data map[string]interface{}) []byte {
	t.Helper()
	body, err := json.Marshal(map[string]interface{}{
		"data": map[string]interface{}{"data": data},
	})
	if err != nil {
		t.Fatal(err)
	}
	return body
}

func TestReadSecret_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != "test-token" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(kvV2Response(t, map[string]interface{}{"DB_PASS": "secret123"}))
	}))
	defer ts.Close()

	c := vault.NewClient(ts.URL, "test-token")
	secrets, err := c.ReadSecret("secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secrets["DB_PASS"] != "secret123" {
		t.Errorf("expected secret123, got %q", secrets["DB_PASS"])
	}
}

func TestReadSecret_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":["no secret"]}`)) //nolint
	}))
	defer ts.Close()

	c := vault.NewClient(ts.URL, "token")
	_, err := c.ReadSecret("secret/data/missing")
	if err == nil {
		t.Fatal("expected error for 404 response")
	}
}

func TestNewClient_FallsBackToEnv(t *testing.T) {
	t.Setenv("VAULT_ADDR", "http://127.0.0.1:8200")
	t.Setenv("VAULT_TOKEN", "env-token")

	c := vault.NewClient("", "")
	if c.Address != "http://127.0.0.1:8200" {
		t.Errorf("unexpected address: %s", c.Address)
	}
	if c.Token != "env-token" {
		t.Errorf("unexpected token: %s", c.Token)
	}
}
