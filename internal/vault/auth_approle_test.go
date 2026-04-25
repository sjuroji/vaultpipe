package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func approleLoginResponse(token string) map[string]interface{} {
	return map[string]interface{}{
		"auth": map[string]interface{}{
			"client_token":   token,
			"lease_duration": 3600,
			"renewable":      true,
		},
	}
}

func TestAppRoleLogin_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/approle/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(approleLoginResponse("s.testtoken123"))
	}))
	defer ts.Close()

	c, _ := NewClient(ts.URL, "")
	token, err := AppRoleLogin(c, "my-role-id", "my-secret-id", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "s.testtoken123" {
		t.Errorf("expected s.testtoken123, got %s", token)
	}
}

func TestAppRoleLogin_CustomMount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/custom-approle/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(approleLoginResponse("s.customtoken"))
	}))
	defer ts.Close()

	c, _ := NewClient(ts.URL, "")
	token, err := AppRoleLogin(c, "role", "secret", "custom-approle")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "s.customtoken" {
		t.Errorf("expected s.customtoken, got %s", token)
	}
}

func TestAppRoleLogin_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c, _ := NewClient(ts.URL, "")
	_, err := AppRoleLogin(c, "bad-role", "bad-secret", "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAppRoleLogin_EmptyToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(approleLoginResponse(""))
	}))
	defer ts.Close()

	c, _ := NewClient(ts.URL, "")
	_, err := AppRoleLogin(c, "role", "secret", "")
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}
