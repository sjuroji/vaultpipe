package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func wrapResponse(token, creationPath string, ttl int) interface{} {
	return map[string]interface{}{
		"wrap_info": map[string]interface{}{
			"token":         token,
			"ttl":           ttl,
			"creation_time": "2024-01-15T10:00:00Z",
			"creation_path": creationPath,
		},
	}
}

func TestWrapSecret_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Wrap-TTL") == "" {
			t.Error("expected X-Vault-Wrap-TTL header")
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(wrapResponse("wrapping-token-abc", "secret/data/myapp", 30))
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "root", HTTP: ts.Client()}
	ws, err := WrapSecret(c, "secret/data/myapp", "30s")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ws.Token != "wrapping-token-abc" {
		t.Errorf("expected token 'wrapping-token-abc', got %q", ws.Token)
	}
	if ws.TTL.Seconds() != 30 {
		t.Errorf("expected TTL 30s, got %v", ws.TTL)
	}
	if ws.CreationPath != "secret/data/myapp" {
		t.Errorf("unexpected creation path: %q", ws.CreationPath)
	}
}

func TestWrapSecret_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "bad-token", HTTP: ts.Client()}
	_, err := WrapSecret(c, "secret/data/myapp", "30s")
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestUnwrapSecret_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/wrapping/unwrap" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"username": "admin",
				"password": "s3cr3t",
			},
		})
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "root", HTTP: ts.Client()}
	data, err := UnwrapSecret(c, "wrapping-token-abc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["username"] != "admin" {
		t.Errorf("expected username 'admin', got %q", data["username"])
	}
	if data["password"] != "s3cr3t" {
		t.Errorf("expected password 's3cr3t', got %q", data["password"])
	}
}

func TestUnwrapSecret_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "root", HTTP: ts.Client()}
	_, err := UnwrapSecret(c, "expired-token")
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}
