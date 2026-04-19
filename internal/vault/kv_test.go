package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func kvDataResponse(data map[string]interface{}) []byte {
	body, _ := json.Marshal(map[string]interface{}{
		"data": map[string]interface{}{
			"data": data,
		},
	})
	return body
}

func TestReadKV_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != "test-token" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(kvDataResponse(map[string]interface{}{"DB_PASS": "secret123", "API_KEY": "abc"}))
	}))
	defer ts.Close()

	c, _ := NewClient(ts.URL, "test-token")
	data, err := c.ReadKV("secret", "myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["DB_PASS"] != "secret123" {
		t.Errorf("expected DB_PASS=secret123, got %q", data["DB_PASS"])
	}
	if data["API_KEY"] != "abc" {
		t.Errorf("expected API_KEY=abc, got %q", data["API_KEY"])
	}
}

func TestReadKV_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errors":["secret not found"]}`))
	}))
	defer ts.Close()

	c, _ := NewClient(ts.URL, "token")
	_, err := c.ReadKV("secret", "missing")
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestReadKV_NonStringValue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(kvDataResponse(map[string]interface{}{"COUNT": 42}))
	}))
	defer ts.Close()

	c, _ := NewClient(ts.URL, "token")
	data, err := c.ReadKV("secret", "myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if data["COUNT"] != "42" {
		t.Errorf("expected COUNT=42, got %q", data["COUNT"])
	}
}
