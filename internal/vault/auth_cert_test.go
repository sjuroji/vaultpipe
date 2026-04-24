package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func certResponse(token string, renewable bool) map[string]interface{} {
	return map[string]interface{}{
		"auth": map[string]interface{}{
			"client_token":   token,
			"lease_duration": 3600,
			"renewable":      renewable,
		},
	}
}

func TestCertLogin_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/cert/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(certResponse("s.certtoken", true))
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	resp, err := CertLogin(c, "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Token != "s.certtoken" {
		t.Errorf("expected token s.certtoken, got %s", resp.Token)
	}
	if !resp.Renewable {
		t.Error("expected renewable to be true")
	}
	if resp.LeaseDur != 3600 {
		t.Errorf("expected lease duration 3600, got %d", resp.LeaseDur)
	}
}

func TestCertLogin_CustomMount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/pki/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(certResponse("s.pkitoken", false))
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	resp, err := CertLogin(c, "pki", "my-role")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Token != "s.pkitoken" {
		t.Errorf("expected s.pkitoken, got %s", resp.Token)
	}
}

func TestCertLogin_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	_, err := CertLogin(c, "", "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCertLogin_EmptyToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(certResponse("", false))
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	_, err := CertLogin(c, "", "")
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}
