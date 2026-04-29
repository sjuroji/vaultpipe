package vault_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yourusername/vaultpipe/internal/vault"
)

func cfResponse(token string) map[string]interface{} {
	return map[string]interface{}{
		"auth": map[string]interface{}{
			"client_token": token,
		},
	}
}

func TestCFLogin_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/cf/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cfResponse("cf-token-abc"))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	tok, err := vault.CFLogin(c, "my-role", "cert-contents", "key-contents", "2024-01-01T00:00:00Z", "sig", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok != "cf-token-abc" {
		t.Errorf("expected cf-token-abc, got %s", tok)
	}
}

func TestCFLogin_CustomMount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/cloudfoundry/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cfResponse("cf-token-xyz"))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	tok, err := vault.CFLogin(c, "my-role", "cert", "key", "time", "sig", "cloudfoundry")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok != "cf-token-xyz" {
		t.Errorf("expected cf-token-xyz, got %s", tok)
	}
}

func TestCFLogin_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	_, err := vault.CFLogin(c, "role", "cert", "key", "time", "sig", "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCFLogin_EmptyToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cfResponse(""))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	_, err := vault.CFLogin(c, "role", "cert", "key", "time", "sig", "")
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}
