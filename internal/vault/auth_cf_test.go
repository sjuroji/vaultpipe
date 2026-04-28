package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
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
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cfResponse("cf-token-abc"))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	tok, err := CFLogin(c, "", "my-role", "2024-01-01T00:00:00Z", "cert-data", "sig-data")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok != "cf-token-abc" {
		t.Errorf("expected cf-token-abc, got %s", tok)
	}
}

func TestCFLogin_CustomMount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/pcf/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cfResponse("pcf-token-xyz"))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	tok, err := CFLogin(c, "pcf", "role", "time", "cert", "sig")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok != "pcf-token-xyz" {
		t.Errorf("expected pcf-token-xyz, got %s", tok)
	}
}

func TestCFLogin_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	_, err := CFLogin(c, "", "role", "time", "cert", "sig")
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestCFLogin_EmptyToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cfResponse(""))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	_, err := CFLogin(c, "", "role", "time", "cert", "sig")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}
