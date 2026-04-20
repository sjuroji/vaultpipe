package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func sealStatusResponse(sealed bool, initialized bool) map[string]interface{} {
	return map[string]interface{}{
		"sealed":       sealed,
		"initialized":  initialized,
		"t":            3,
		"n":            5,
		"progress":     0,
		"version":      "1.13.0",
		"cluster_name": "vault-cluster",
		"cluster_id":   "abc-123",
	}
}

func TestGetSealStatus_Unsealed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/seal-status" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(sealStatusResponse(false, true))
	}))
	defer ts.Close()

	status, err := GetSealStatus(ts.Client(), ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Sealed {
		t.Error("expected sealed=false")
	}
	if !status.Initialized {
		t.Error("expected initialized=true")
	}
	if status.Version != "1.13.0" {
		t.Errorf("unexpected version: %s", status.Version)
	}
}

func TestGetSealStatus_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()

	_, err := GetSealStatus(ts.Client(), ts.URL)
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestIsSealed_True(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(sealStatusResponse(true, true))
	}))
	defer ts.Close()

	sealed, err := IsSealed(ts.Client(), ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !sealed {
		t.Error("expected sealed=true")
	}
}

func TestIsSealed_False(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(sealStatusResponse(false, true))
	}))
	defer ts.Close()

	sealed, err := IsSealed(ts.Client(), ts.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sealed {
		t.Error("expected sealed=false")
	}
}
