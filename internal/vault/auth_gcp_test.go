package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func gcpResponse(token string) map[string]interface{} {
	return map[string]interface{}{
		"auth": map[string]interface{}{
			"client_token": token,
		},
	}
}

func TestGCPLogin_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/gcp/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(gcpResponse("gcp-token-abc"))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	resp, err := GCPLogin(c, GCPLoginRequest{Role: "my-role", JWT: "signed-jwt"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ClientToken != "gcp-token-abc" {
		t.Errorf("expected gcp-token-abc, got %s", resp.ClientToken)
	}
}

func TestGCPLogin_CustomMount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/gcp-prod/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(gcpResponse("gcp-prod-token"))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	resp, err := GCPLogin(c, GCPLoginRequest{Role: "my-role", JWT: "signed-jwt", Mount: "gcp-prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ClientToken != "gcp-prod-token" {
		t.Errorf("expected gcp-prod-token, got %s", resp.ClientToken)
	}
}

func TestGCPLogin_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	_, err := GCPLogin(c, GCPLoginRequest{Role: "my-role", JWT: "bad-jwt"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGCPLogin_EmptyToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(gcpResponse(""))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	_, err := GCPLogin(c, GCPLoginRequest{Role: "my-role", JWT: "signed-jwt"})
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}
