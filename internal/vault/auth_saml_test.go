package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func samlResponse(token string) map[string]interface{} {
	return map[string]interface{}{
		"auth": map[string]interface{}{
			"client_token":   token,
			"lease_duration": 3600,
			"renewable":      true,
		},
	}
}

func TestSAMLLogin_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/saml/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(samlResponse("s.samltoken"))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	resp, err := SAMLLogin(c, SAMLLoginRequest{SAMLResponse: "base64encodedsaml=="})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ClientToken != "s.samltoken" {
		t.Errorf("expected s.samltoken, got %s", resp.ClientToken)
	}
	if resp.LeaseDuration != 3600 {
		t.Errorf("expected lease 3600, got %d", resp.LeaseDuration)
	}
	if !resp.Renewable {
		t.Error("expected renewable to be true")
	}
}

func TestSAMLLogin_CustomMount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/corp-saml/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(samlResponse("s.customtoken"))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	resp, err := SAMLLogin(c, SAMLLoginRequest{SAMLResponse: "assertion==", Mount: "corp-saml"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ClientToken != "s.customtoken" {
		t.Errorf("expected s.customtoken, got %s", resp.ClientToken)
	}
}

func TestSAMLLogin_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	_, err := SAMLLogin(c, SAMLLoginRequest{SAMLResponse: "bad=="})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSAMLLogin_EmptyToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(samlResponse(""))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	_, err := SAMLLogin(c, SAMLLoginRequest{SAMLResponse: "assertion=="})
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}
