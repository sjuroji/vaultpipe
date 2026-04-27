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

func TestRunOIDCLogin_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"auth": map[string]interface{}{
				"client_token": "oidc-test-token",
			},
		})
	}))
	defer ts.Close()

	t.Setenv("VAULT_ADDR", ts.URL)
	t.Setenv("VAULT_TOKEN", "root")
	t.Setenv("OIDC_JWT_TOKEN", "my.test.jwt")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"oidc"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(buf.String(), "oidc-test-token") {
		t.Errorf("expected token in output, got: %s", buf.String())
	}
}

func TestRunOIDCLogin_ExportFlag(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"auth": map[string]interface{}{
				"client_token": "oidc-export-token",
			},
		})
	}))
	defer ts.Close()

	t.Setenv("VAULT_ADDR", ts.URL)
	t.Setenv("VAULT_TOKEN", "root")
	t.Setenv("OIDC_JWT_TOKEN", "my.test.jwt")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"oidc", "--export"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()
	if !strings.HasPrefix(out, "export VAULT_TOKEN=") {
		t.Errorf("expected export statement, got: %s", out)
	}
}

func TestRunOIDCLogin_MissingEnvVar(t *testing.T) {
	os.Unsetenv("OIDC_JWT_TOKEN")

	rootCmd.SetArgs([]string{"oidc"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when OIDC_JWT_TOKEN is missing")
	}
}
