package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func capabilitiesResponse(caps []string) map[string]interface{} {
	return map[string]interface{}{
		"capabilities": caps,
	}
}

func TestCheckCapabilities_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/capabilities-self" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(capabilitiesResponse([]string{"read", "list"}))
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "test-token", HTTPClient: ts.Client()}
	caps, err := CheckCapabilities(c, "secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(caps) != 2 {
		t.Fatalf("expected 2 capabilities, got %d", len(caps))
	}
	if caps[0] != "read" {
		t.Errorf("expected first cap to be 'read', got %q", caps[0])
	}
}

func TestCheckCapabilities_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "bad-token", HTTPClient: ts.Client()}
	_, err := CheckCapabilities(c, "secret/data/myapp")
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestHasCapability_True(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(capabilitiesResponse([]string{"read", "list"}))
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "test-token", HTTPClient: ts.Client()}
	ok, err := HasCapability(c, "secret/data/myapp", "read")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Error("expected HasCapability to return true for 'read'")
	}
}

func TestHasCapability_False(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(capabilitiesResponse([]string{"list"}))
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "test-token", HTTPClient: ts.Client()}
	ok, err := HasCapability(c, "secret/data/myapp", "delete")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Error("expected HasCapability to return false for 'delete'")
	}
}
