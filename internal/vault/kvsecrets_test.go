package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func kvListResponse(keys []string) []byte {
	body := map[string]interface{}{
		"data": map[string]interface{}{
			"keys": keys,
		},
	}
	b, _ := json.Marshal(body)
	return b
}

func TestListKV_OK(t *testing.T) {
	expected := []string{"foo", "bar", "baz"}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "LIST" {
			t.Errorf("expected LIST method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(kvListResponse(expected))
	}))
	defer ts.Close()

	client := &Client{Address: ts.URL, Token: "test-token", HTTP: ts.Client()}
	result, err := ListKV(client, "secret", "myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Keys) != len(expected) {
		t.Fatalf("expected %d keys, got %d", len(expected), len(result.Keys))
	}
	for i, k := range expected {
		if result.Keys[i] != k {
			t.Errorf("key[%d]: expected %q, got %q", i, k, result.Keys[i])
		}
	}
}

func TestListKV_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	client := &Client{Address: ts.URL, Token: "bad-token", HTTP: ts.Client()}
	_, err := ListKV(client, "secret", "myapp")
	if err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}

func TestDeleteKV_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE method, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	client := &Client{Address: ts.URL, Token: "test-token", HTTP: ts.Client()}
	if err := DeleteKV(client, "secret", "myapp/config"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteKV_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	client := &Client{Address: ts.URL, Token: "bad-token", HTTP: ts.Client()}
	if err := DeleteKV(client, "secret", "myapp/config"); err == nil {
		t.Fatal("expected error for non-OK status, got nil")
	}
}
