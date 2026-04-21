package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestRunRekeyStatus_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"started":  true,
			"nonce":    "test-nonce",
			"progress": 1,
			"required": 3,
			"t":        3,
			"n":        5,
		})
	}))
	defer ts.Close()

	t.Setenv("VAULT_ADDR", ts.URL)
	t.Setenv("VAULT_TOKEN", "root")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"rekey", "status"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "nonce=test-nonce") {
		t.Errorf("expected nonce in output, got: %s", buf.String())
	}
}

func TestRunRekeyCancel_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	os.Setenv("VAULT_ADDR", ts.URL)
	os.Setenv("VAULT_TOKEN", "root")
	defer os.Unsetenv("VAULT_ADDR")
	defer os.Unsetenv("VAULT_TOKEN")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"rekey", "cancel"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "cancelled") {
		t.Errorf("expected cancelled in output, got: %s", buf.String())
	}
}

func TestRekeyInitFlags_Defaults(t *testing.T) {
	f := rekeyInitCmd.Flags()
	shares, err := f.GetInt("shares")
	if err != nil {
		t.Fatalf("shares flag not found: %v", err)
	}
	if shares != 5 {
		t.Errorf("expected default shares=5, got %d", shares)
	}
	threshold, err := f.GetInt("threshold")
	if err != nil {
		t.Fatalf("threshold flag not found: %v", err)
	}
	if threshold != 3 {
		t.Errorf("expected default threshold=3, got %d", threshold)
	}
}
