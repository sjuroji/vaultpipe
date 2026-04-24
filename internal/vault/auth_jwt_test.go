package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func jwtResponse(token string) map[string]interface{} {
	return map[string]interface{}{
		"auth": map[string]interface{}{
			"client_token": token,
		},
	}
}

func TestJWTLogin_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/auth/jwt/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(jwtResponse("jwt-token-abc"))
	}))
	defer ts.Close()

	client := newTestClient(ts.URL)
	resp, err := JWTLogin(client, "my.jwt.token", "my-role", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Token != "jwt-token-abc" {
		t.Errorf("expected jwt-token-abc, got %s", resp.Token)
	}
}

func TestJWTLogin_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	client := newTestClient(ts.URL)
	_, err := JWTLogin(client, "bad.jwt", "role", "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestJWTLogin_EmptyToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(jwtResponse(""))
	}))
	defer ts.Close()

	client := newTestClient(ts.URL)
	_, err := JWTLogin(client, "some.jwt", "role", "")
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}

func TestJWTLogin_CustomMount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/oidc/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(jwtResponse("oidc-token-xyz"))
	}))
	defer ts.Close()

	client := newTestClient(ts.URL)
	resp, err := JWTLogin(client, "my.jwt", "role", "oidc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Token != "oidc-token-xyz" {
		t.Errorf("expected oidc-token-xyz, got %s", resp.Token)
	}
}
