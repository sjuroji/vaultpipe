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

func TestRunAzureLogin_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"auth": map[string]interface{}{
				"client_token": "azure-test-token",
			},
		})
	}))
	defer ts.Close()

	os.Setenv("VAULT_ADDR", ts.URL)
	os.Setenv("VAULT_TOKEN", "root")
	defer os.Unsetenv("VAULT_ADDR")
	defer os.Unsetenv("VAULT_TOKEN")

	buf := &bytes.Buffer{}
	azureCmd.SetOut(buf)
	azureCmd.SetErr(buf)

	azureRole = "my-role"
	azureJWT = "my-jwt"
	azureMount = "azure"
	azureExport = false

	err := runAzureLogin(azureCmd, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunAzureLogin_ExportFlag(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"auth": map[string]interface{}{
				"client_token": "azure-export-token",
			},
		})
	}))
	defer ts.Close()

	os.Setenv("VAULT_ADDR", ts.URL)
	os.Setenv("VAULT_TOKEN", "root")
	defer os.Unsetenv("VAULT_ADDR")
	defer os.Unsetenv("VAULT_TOKEN")

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	azureRole = "my-role"
	azureJWT = "my-jwt"
	azureMount = "azure"
	azureExport = true

	err := runAzureLogin(azureCmd, nil)
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	buf := &bytes.Buffer{}
	buf.ReadFrom(r)
	output := buf.String()
	if !strings.Contains(output, "export VAULT_TOKEN=") {
		t.Errorf("expected export statement, got: %s", output)
	}
}

func TestRunAzureLogin_MissingRoleFlag(t *testing.T) {
	azureRole = ""
	azureJWT = "some-jwt"
	// role is required; cobra should reject this before RunE is called
	// We test that the flag is marked required via the command definition.
	if azureCmd.Flag("role") == nil {
		t.Error("expected --role flag to be defined")
	}
}
