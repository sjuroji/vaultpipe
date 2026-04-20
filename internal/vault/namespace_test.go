package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func namespacesResponse(keys []string) []byte {
	payload := map[string]any{
		"data": map[string]any{
			"keys": keys,
		},
	}
	b, _ := json.Marshal(payload)
	return b
}

func TestListNamespaces_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != "test-token" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(namespacesResponse([]string{"ns1/", "ns2/"}))
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "test-token", HTTP: ts.Client()}
	keys, err := ListNamespaces(c, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 2 || keys[0] != "ns1/" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestListNamespaces_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("permission denied"))
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "bad-token", HTTP: ts.Client()}
	_, err := ListNamespaces(c, "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetNamespace_OK(t *testing.T) {
	expected := NamespaceInfo{Path: "team-a/", ID: "abc123", Meta: map[string]string{"owner": "alice"}}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		b, _ := json.Marshal(expected)
		w.Write(b)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "test-token", HTTP: ts.Client()}
	info, err := GetNamespace(c, "team-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.ID != "abc123" {
		t.Errorf("expected ID abc123, got %s", info.ID)
	}
	if info.Meta["owner"] != "alice" {
		t.Errorf("expected owner alice, got %s", info.Meta["owner"])
	}
}

func TestGetNamespace_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("namespace not found"))
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "test-token", HTTP: ts.Client()}
	_, err := GetNamespace(c, "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
