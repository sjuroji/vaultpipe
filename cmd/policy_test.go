package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRunPolicyList_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"keys": []string{"default", "admin"},
			},
		})
	}))
	defer ts.Close()

	t.Setenv("VAULT_ADDR", ts.URL)
	t.Setenv("VAULT_TOKEN", "test-token")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"policy", "list"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunPolicyGet_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"name":  "admin",
				"rules": `path "secret/*" { capabilities = ["read"] }`,
			},
		})
	}))
	defer ts.Close()

	t.Setenv("VAULT_ADDR", ts.URL)
	t.Setenv("VAULT_TOKEN", "test-token")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"policy", "get", "admin"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "admin") {
		t.Errorf("expected output to contain policy name, got: %s", output)
	}
}

func TestRunPolicyGet_MissingArg(t *testing.T) {
	rootCmd.SetArgs([]string{"policy", "get"})
	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing argument")
	}
}
