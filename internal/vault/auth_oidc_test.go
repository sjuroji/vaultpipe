package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func oidcResponse(token string) map[string]interface{} {
	return map[string]interface{}{
		"auth": map[string]interface{}{
			"client_token": token,
		},
	}
}

func TestOIDCLogin_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/oidc/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(oidcResponse("oidc-token-abc"))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	resp, err := OIDCLogin(c, "my.jwt.token", "myrole", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Token != "oidc-token-abc" {
		t.Errorf("expected token %q, got %q", "oidc-token-abc", resp.Token)
	}
}

func TestOIDCLogin_CustomMount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/myoidc/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(oidcResponse("custom-oidc-token"))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	resp, err := OIDCLogin(c, "my.jwt.token", "myrole", "myoidc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Token != "custom-oidc-token" {
		t.Errorf("expected token %q, got %q", "custom-oidc-token", resp.Token)
	}
}

func TestOIDCLogin_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	_, err := OIDCLogin(c, "bad.jwt", "role", "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOIDCLogin_EmptyToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(oidcResponse(""))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	_, err := OIDCLogin(c, "my.jwt.token", "role", "")
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}
