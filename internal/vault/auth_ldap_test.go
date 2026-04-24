package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func ldapResponse(token string, policies []string) map[string]interface{} {
	return map[string]interface{}{
		"auth": map[string]interface{}{
			"client_token":   token,
			"policies":       policies,
			"lease_duration": 3600,
		},
	}
}

func TestLDAPLogin_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/auth/ldap/login/alice" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ldapResponse("s.ldaptoken", []string{"default"}))
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	resp, err := LDAPLogin(c, "alice", "secret", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Token != "s.ldaptoken" {
		t.Errorf("expected s.ldaptoken, got %s", resp.Token)
	}
	if len(resp.Policies) != 1 || resp.Policies[0] != "default" {
		t.Errorf("unexpected policies: %v", resp.Policies)
	}
	if resp.LeaseDur != 3600 {
		t.Errorf("expected lease 3600, got %d", resp.LeaseDur)
	}
}

func TestLDAPLogin_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	_, err := LDAPLogin(c, "bob", "wrong", "ldap")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLDAPLogin_EmptyToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ldapResponse("", nil))
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	_, err := LDAPLogin(c, "alice", "secret", "ldap")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestLDAPLogin_CustomMount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/corp-ldap/login/alice" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ldapResponse("s.custom", nil))
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	resp, err := LDAPLogin(c, "alice", "secret", "corp-ldap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Token != "s.custom" {
		t.Errorf("expected s.custom, got %s", resp.Token)
	}
}
