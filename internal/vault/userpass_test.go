package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func userpassResponse(token string, renewable bool) map[string]interface{} {
	return map[string]interface{}{
		"auth": map[string]interface{}{
			"client_token":   token,
			"lease_duration": 3600,
			"renewable":      renewable,
		},
	}
}

func TestUserpassLogin_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/userpass/login/alice" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userpassResponse("tok-abc", true))
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, HTTP: ts.Client()}
	resp, err := UserpassLogin(c, "alice", "s3cr3t", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Token != "tok-abc" {
		t.Errorf("expected tok-abc, got %s", resp.Token)
	}
	if !resp.Renewable {
		t.Error("expected renewable=true")
	}
	if resp.LeaseDur != 3600 {
		t.Errorf("expected lease 3600, got %d", resp.LeaseDur)
	}
}

func TestUserpassLogin_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, HTTP: ts.Client()}
	_, err := UserpassLogin(c, "bob", "wrong", "userpass")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUserpassLogin_EmptyToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userpassResponse("", false))
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, HTTP: ts.Client()}
	_, err := UserpassLogin(c, "alice", "pass", "")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestUserpassLogin_CustomMount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/custom-mount/login/alice" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(userpassResponse("tok-xyz", false))
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, HTTP: ts.Client()}
	resp, err := UserpassLogin(c, "alice", "pass", "custom-mount")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Token != "tok-xyz" {
		t.Errorf("expected tok-xyz, got %s", resp.Token)
	}
}
