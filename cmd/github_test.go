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

func TestRunGitHubLogin_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"auth": map[string]interface{}{
				"client_token":   "gh-vault-token",
				"lease_duration": 3600,
				"renewable":      true,
			},
		})
	}))
	defer ts.Close()

	t.Setenv("VAULT_ADDR", ts.URL)
	t.Setenv("VAULT_TOKEN", "root")
	t.Setenv("GITHUB_TOKEN", "mypat")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"github"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "gh-vault-token") {
		t.Errorf("expected token in output, got: %s", buf.String())
	}
}

func TestRunGitHubLogin_ExportFlag(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"auth": map[string]interface{}{
				"client_token":   "export-token",
				"lease_duration": 1800,
				"renewable":      false,
			},
		})
	}))
	defer ts.Close()

	t.Setenv("VAULT_ADDR", ts.URL)
	t.Setenv("VAULT_TOKEN", "root")
	t.Setenv("GITHUB_TOKEN", "mypat")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"github", "--export"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "export VAULT_TOKEN=export-token") {
		t.Errorf("expected export statement, got: %s", buf.String())
	}
}

func TestRunGitHubLogin_MissingEnvVar(t *testing.T) {
	os.Unsetenv("GITHUB_TOKEN")

	rootCmd.SetArgs([]string{"github"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when GITHUB_TOKEN is missing")
	}
}
