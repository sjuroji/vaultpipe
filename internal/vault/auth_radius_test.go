package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func radiusResponse(token string) map[string]interface{} {
	return map[string]interface{}{
		"auth": map[string]interface{}{
			"client_token": token,
		},
	}
}

func TestRADIUSLogin_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(radiusResponse("radius-token-abc"))
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	token, err := RADIUSLogin(c, "alice", "secret", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "radius-token-abc" {
		t.Errorf("expected radius-token-abc, got %s", token)
	}
}

func TestRADIUSLogin_CustomMount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/corp-radius/login/bob" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(radiusResponse("custom-mount-token"))
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	token, err := RADIUSLogin(c, "bob", "pass", "corp-radius")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "custom-mount-token" {
		t.Errorf("expected custom-mount-token, got %s", token)
	}
}

func TestRADIUSLogin_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	_, err := RADIUSLogin(c, "alice", "wrong", "")
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestRADIUSLogin_EmptyToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(radiusResponse(""))
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	_, err := RADIUSLogin(c, "alice", "secret", "")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}
