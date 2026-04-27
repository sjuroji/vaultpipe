package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func tlsLoginResponse(token string) map[string]interface{} {
	return map[string]interface{}{
		"auth": map[string]interface{}{
			"client_token": token,
		},
	}
}

func TestTLSLogin_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/auth/cert/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tlsLoginResponse("tls-token-abc"))
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	resp, err := TLSLogin(c, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Token != "tls-token-abc" {
		t.Errorf("expected token tls-token-abc, got %s", resp.Token)
	}
	if resp.Mount != "cert" {
		t.Errorf("expected mount cert, got %s", resp.Mount)
	}
}

func TestTLSLogin_CustomMount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/pki-tls/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tlsLoginResponse("tls-custom-token"))
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	resp, err := TLSLogin(c, "pki-tls")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Mount != "pki-tls" {
		t.Errorf("expected mount pki-tls, got %s", resp.Mount)
	}
}

func TestTLSLogin_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	_, err := TLSLogin(c, "")
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestTLSLogin_EmptyToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tlsLoginResponse(""))
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	_, err := TLSLogin(c, "")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}
