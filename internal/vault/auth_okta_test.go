package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func oktaResponse(token string) map[string]interface{} {
	return map[string]interface{}{
		"auth": map[string]interface{}{
			"client_token": token,
		},
	}
}

func TestOktaLogin_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/auth/okta/login/alice" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(oktaResponse("s.oktatoken"))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	token, err := OktaLogin(c, "alice", "secret", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "s.oktatoken" {
		t.Errorf("expected s.oktatoken, got %s", token)
	}
}

func TestOktaLogin_CustomMount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/corp-okta/login/bob" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(oktaResponse("s.customtoken"))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	token, err := OktaLogin(c, "bob", "pass", "corp-okta")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "s.customtoken" {
		t.Errorf("expected s.customtoken, got %s", token)
	}
}

func TestOktaLogin_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	_, err := OktaLogin(c, "alice", "wrong", "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOktaLogin_EmptyToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(oktaResponse(""))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	_, err := OktaLogin(c, "alice", "secret", "")
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}
