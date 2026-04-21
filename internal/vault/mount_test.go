package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func mountsResponse(t *testing.T, w http.ResponseWriter, data MountsResponse, status int) {
	t.Helper()
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		t.Fatalf("encode mounts response: %v", err)
	}
}

func TestListMounts_OK(t *testing.T) {
	expected := MountsResponse{
		"secret/": {Type: "kv", Description: "key/value store", Accessor: "kv_abc123"},
		"sys/":    {Type: "system", Description: "system backend", Accessor: "sys_xyz"},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/mounts" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		mountsResponse(t, w, expected, http.StatusOK)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "test-token", HTTP: ts.Client()}
	got, err := ListMounts(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["secret/"].Type != "kv" {
		t.Errorf("expected kv type, got %q", got["secret/"].Type)
	}
	if len(got) != 2 {
		t.Errorf("expected 2 mounts, got %d", len(got))
	}
}

func TestListMounts_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "bad-token", HTTP: ts.Client()}
	_, err := ListMounts(c)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetMount_OK(t *testing.T) {
	data := MountsResponse{
		"secret/": {Type: "kv", Description: "kv store"},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mountsResponse(t, w, data, http.StatusOK)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "tok", HTTP: ts.Client()}
	info, err := GetMount(c, "secret/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Type != "kv" {
		t.Errorf("expected type kv, got %q", info.Type)
	}
}

func TestGetMount_NotFound(t *testing.T) {
	data := MountsResponse{
		"secret/": {Type: "kv"},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mountsResponse(t, w, data, http.StatusOK)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "tok", HTTP: ts.Client()}
	_, err := GetMount(c, "pki/")
	if err == nil {
		t.Fatal("expected error for missing mount, got nil")
	}
}
