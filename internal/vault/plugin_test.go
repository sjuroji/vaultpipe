package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func pluginCatalogResponse() map[string]interface{} {
	return map[string]interface{}{
		"data": map[string]interface{}{
			"auth":     []string{"approle", "github"},
			"database": []string{"mysql", "postgresql"},
			"secret":   []string{"kv", "pki"},
		},
	}
}

func TestListPlugins_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/plugins/catalog" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(pluginCatalogResponse())
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "test-token", HTTP: ts.Client()}
	result, err := ListPlugins(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Data.Auth) != 2 {
		t.Errorf("expected 2 auth plugins, got %d", len(result.Data.Auth))
	}
	if len(result.Data.Secret) != 2 {
		t.Errorf("expected 2 secret plugins, got %d", len(result.Data.Secret))
	}
}

func TestListPlugins_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "bad-token", HTTP: ts.Client()}
	_, err := ListPlugins(c)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetPlugin_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"name":    "kv",
				"type":    "secret",
				"version": "v2",
				"builtin": true,
			},
		})
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "test-token", HTTP: ts.Client()}
	plugin, err := GetPlugin(c, "secret", "kv")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if plugin.Name != "kv" {
		t.Errorf("expected name 'kv', got %q", plugin.Name)
	}
	if !plugin.Builtin {
		t.Error("expected builtin to be true")
	}
}

func TestGetPlugin_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "test-token", HTTP: ts.Client()}
	_, err := GetPlugin(c, "secret", "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
