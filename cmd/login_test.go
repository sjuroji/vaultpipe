package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRunLogin_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"auth": map[string]interface{}{
				"client_token":   "tok-login-ok",
				"lease_duration": 7200,
				"renewable":      true,
			},
		})
	}))
	defer ts.Close()
	t.Setenv("VAULT_ADDR", ts.URL)
	t.Setenv("VAULT_TOKEN", "root")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"login", "-u", "alice", "-p", "secret"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "tok-login-ok") {
		t.Errorf("expected token in output, got: %s", buf.String())
	}
}

func TestRunLogin_ExportFlag(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"auth": map[string]interface{}{
				"client_token":   "tok-export",
				"lease_duration": 3600,
				"renewable":      false,
			},
		})
	}))
	defer ts.Close()
	t.Setenv("VAULT_ADDR", ts.URL)
	t.Setenv("VAULT_TOKEN", "root")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"login", "-u", "alice", "-p", "secret", "--export"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "export VAULT_TOKEN=tok-export") {
		t.Errorf("expected export statement, got: %s", out)
	}
}
