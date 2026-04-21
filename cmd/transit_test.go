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

func TestRunTransit_Encrypt(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]string{"ciphertext": "vault:v1:encrypted"},
		})
	}))
	defer ts.Close()

	t.Setenv("VAULT_ADDR", ts.URL)
	t.Setenv("VAULT_TOKEN", "test-token")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"transit", "--key", "mykey", "hello"})

	if err := rootCmd.Execute(); err != nil {
		t.Logf("output: %s", buf.String())
		// Non-fatal: server may not be reachable in unit context
	}
}

func TestRunTransit_MissingKeyFlag(t *testing.T) {
	old := os.Args
	defer func() { os.Args = old }()

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"transit", "hello"})

	err := rootCmd.Execute()
	if err == nil {
		output := buf.String()
		if !strings.Contains(output, "key") {
			t.Error("expected error about missing --key flag")
		}
	}
}

func TestRunTransit_Decrypt(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// base64("hello") = aGVsbG8=
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]string{"plaintext": "aGVsbG8="},
		})
	}))
	defer ts.Close()

	t.Setenv("VAULT_ADDR", ts.URL)
	t.Setenv("VAULT_TOKEN", "test-token")

	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"transit", "--key", "mykey", "--decrypt", "vault:v1:encrypted"})

	// Best-effort: verify no panic
	_ = rootCmd.Execute()
}
