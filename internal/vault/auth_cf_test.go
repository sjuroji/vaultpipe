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
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cfResponse("cf-token-abc"))
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	token, err := CFLogin(c, CFLoginRequest{
		RoleID:         "my-role",
		SigningTime:    "2024-01-01T00:00:00Z",
		CFInstanceCert: "cert-data",
		Signature:      "sig-data",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "cf-token-abc" {
		t.Errorf("expected cf-token-abc, got %s", token)
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

	c := newTestClient(t, ts.URL)
	token, err := CFLogin(c, CFLoginRequest{
		RoleID: "my-role",
		Mount:  "cloudfoundry",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "cf-token-xyz" {
		t.Errorf("expected cf-token-xyz, got %s", token)
	}
}

func TestCFLogin_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	_, err := CFLogin(c, CFLoginRequest{RoleID: "bad-role"})
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

	c := newTestClient(t, ts.URL)
	_, err := CFLogin(c, CFLoginRequest{RoleID: "my-role"})
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}
