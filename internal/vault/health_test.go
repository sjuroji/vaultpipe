package vault

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func healthResponse(t *testing.T, status HealthStatus, code int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(status)
	}))
}

func TestCheckHealth_OK(t *testing.T) {
	expected := HealthStatus{
		Initialized: true,
		Sealed:      false,
		Standby:     false,
		Version:     "1.15.0",
		ClusterName: "vault-cluster",
	}
	srv := healthResponse(t, expected, http.StatusOK)
	defer srv.Close()

	got, err := CheckHealth(context.Background(), srv.Client(), srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Version != expected.Version {
		t.Errorf("version: got %q, want %q", got.Version, expected.Version)
	}
	if got.Sealed != false {
		t.Errorf("expected unsealed")
	}
}

func TestCheckHealth_SealedReturnsBody(t *testing.T) {
	sealed := HealthStatus{Initialized: true, Sealed: true}
	// Vault returns 503 when sealed
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(sealed)
	}))
	defer srv.Close()

	_, err := CheckHealth(context.Background(), srv.Client(), srv.URL)
	if err == nil {
		t.Fatal("expected error for 503 status")
	}
}

func TestCheckHealth_StandbyReadable(t *testing.T) {
	standby := HealthStatus{Initialized: true, Sealed: false, Standby: true, Version: "1.15.0"}
	srv := healthResponse(t, standby, 429)
	defer srv.Close()

	got, err := CheckHealth(context.Background(), srv.Client(), srv.URL)
	if err != nil {
		t.Fatalf("unexpected error for standby: %v", err)
	}
	if !got.Standby {
		t.Errorf("expected standby=true")
	}
}

func TestIsHealthy_True(t *testing.T) {
	s := &HealthStatus{Initialized: true, Sealed: false, Standby: false}
	if !IsHealthy(s) {
		t.Error("expected healthy")
	}
}

func TestIsHealthy_False_Sealed(t *testing.T) {
	s := &HealthStatus{Initialized: true, Sealed: true, Standby: false}
	if IsHealthy(s) {
		t.Error("expected unhealthy when sealed")
	}
}

func TestIsHealthy_False_Standby(t *testing.T) {
	s := &HealthStatus{Initialized: true, Sealed: false, Standby: true}
	if IsHealthy(s) {
		t.Error("expected unhealthy when standby")
	}
}
