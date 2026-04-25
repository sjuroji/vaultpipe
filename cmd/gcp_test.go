package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRunGCPLogin_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"auth": map[string]interface{}{
				"client_token": "gcp-test-token",
			},
		})
	}))
	defer ts.Close()

	t.Setenv("VAULT_ADDR", ts.URL)
	t.Setenv("VAULT_TOKEN", "root")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"gcp", "--role", "my-role", "--jwt", "signed-jwt"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "gcp-test-token") {
		t.Errorf("expected token in output, got: %s", buf.String())
	}
}

func TestRunGCPLogin_ExportFlag(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"auth": map[string]interface{}{
				"client_token": "gcp-export-token",
			},
		})
	}))
	defer ts.Close()

	t.Setenv("VAULT_ADDR", ts.URL)
	t.Setenv("VAULT_TOKEN", "root")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"gcp", "--role", "my-role", "--jwt", "signed-jwt", "--export"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "export VAULT_TOKEN=gcp-export-token") {
		t.Errorf("expected export statement in output, got: %s", buf.String())
	}
}

func TestRunGCPLogin_MissingRoleFlag(t *testing.T) {
	rootCmd.SetArgs([]string{"gcp", "--jwt", "some-jwt"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing --role flag, got nil")
	}
}
