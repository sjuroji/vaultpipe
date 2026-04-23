package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func initStatusResponse(initialized bool) map[string]interface{} {
	return map[string]interface{}{
		"initialized": initialized,
	}
}

func TestGetInitStatus_Initialized(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/init" || r.Method != http.MethodGet {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(initStatusResponse(true))
	}))
	defer ts.Close()

	client, _ := NewClient(ts.URL, "test-token")
	status, err := GetInitStatus(client)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.Initialized {
		t.Errorf("expected initialized=true, got false")
	}
}

func TestGetInitStatus_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	client, _ := NewClient(ts.URL, "test-token")
	_, err := GetInitStatus(client)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestInitialize_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"root_token": "s.abc123",
			"keys":       []string{"key1", "key2", "key3"},
		})
	}))
	defer ts.Close()

	client, _ := NewClient(ts.URL, "test-token")
	resp, err := Initialize(client, 3, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.RootToken != "s.abc123" {
		t.Errorf("expected root_token s.abc123, got %s", resp.RootToken)
	}
	if len(resp.UnsealKeys) != 3 {
		t.Errorf("expected 3 unseal keys, got %d", len(resp.UnsealKeys))
	}
}

func TestInitialize_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	client, _ := NewClient(ts.URL, "test-token")
	_, err := Initialize(client, 5, 3)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
