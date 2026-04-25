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

func TestRunAWSLogin_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"auth": map[string]interface{}{"client_token": "aws-tok-123"},
		})
	}))
	defer ts.Close()

	t.Setenv("VAULT_ADDR", ts.URL)
	t.Setenv("VAULT_TOKEN", "test-root")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"aws", "--role", "my-role", "--iam-request-url", "dXJs", "--iam-request-body", "Ym9keQ==", "--iam-request-headers", ""})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunAWSLogin_ExportFlag(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"auth": map[string]interface{}{"client_token": "export-tok"},
		})
	}))
	defer ts.Close()

	t.Setenv("VAULT_ADDR", ts.URL)
	t.Setenv("VAULT_TOKEN", "test-root")

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs([]string{"aws", "--role", "my-role", "--iam-request-url", "dXJs", "--iam-request-body", "Ym9keQ==", "--iam-request-headers", "", "--export"})
	_ = rootCmd.Execute()

	w.Close()
	os.Stdout = old

	var out bytes.Buffer
	out.ReadFrom(r)
	if !strings.Contains(out.String(), "export VAULT_TOKEN=") {
		t.Errorf("expected export statement, got: %s", out.String())
	}
}

func TestRunAWSLogin_MissingRoleFlag(t *testing.T) {
	rootCmd.SetArgs([]string{"aws"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error when --role is missing")
	}
}
