package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func enginesResponse(t *testing.T, engines map[string]SecretEngine) http.HandlerFunc {
	t.Helper()
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(engines)
	}
}

func TestListSecretEngines_OK(t *testing.T) {
	payload := map[string]SecretEngine{
		"secret/": {Type: "kv", Description: "key/value", Accessor: "kv_abc123"},
		"pki/":    {Type: "pki", Description: "PKI", Accessor: "pki_def456"},
	}
	ts := httptest.NewServer(enginesResponse(t, payload))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "test-token", HTTP: ts.Client()}
	engines, err := ListSecretEngines(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(engines) != 2 {
		t.Fatalf("expected 2 engines, got %d", len(engines))
	}
	if engines["secret/"].Type != "kv" {
		t.Errorf("expected type kv, got %s", engines["secret/"].Type)
	}
}

func TestListSecretEngines_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "bad-token", HTTP: ts.Client()}
	_, err := ListSecretEngines(c)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetSecretEngine_OK(t *testing.T) {
	payload := map[string]SecretEngine{
		"secret/": {Type: "kv", Accessor: "kv_abc123"},
	}
	ts := httptest.NewServer(enginesResponse(t, payload))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "test-token", HTTP: ts.Client()}
	engine, err := GetSecretEngine(c, "secret")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if engine.Type != "kv" {
		t.Errorf("expected type kv, got %s", engine.Type)
	}
}

func TestGetSecretEngine_NotFound(t *testing.T) {
	payload := map[string]SecretEngine{
		"secret/": {Type: "kv"},
	}
	ts := httptest.NewServer(enginesResponse(t, payload))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "test-token", HTTP: ts.Client()}
	_, err := GetSecretEngine(c, "nonexistent")
	if err == nil {
		t.Fatal("expected error for missing engine, got nil")
	}
}
