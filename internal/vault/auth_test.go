package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func approleResponse(token string) map[string]interface{} {
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
		json.NewEncoder(w).Encode(approleResponse("s.testtoken"))
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTP: ts.Client()}
	tok, err := AppRoleLogin(c, "my-role-id", "my-secret-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok != "s.testtoken" {
		t.Errorf("expected s.testtoken, got %s", tok)
	}
}

func TestAppRoleLogin_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTP: ts.Client()}
	_, err := AppRoleLogin(c, "bad-role", "bad-secret")
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestAppRoleLogin_EmptyToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(approleResponse(""))
	}))
	defer ts.Close()

	c := &Client{BaseURL: ts.URL, HTTP: ts.Client()}
	_, err := AppRoleLogin(c, "role", "secret")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}
