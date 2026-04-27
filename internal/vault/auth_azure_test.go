package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func azureResponse(token string) map[string]interface{} {
	return map[string]interface{}{
		"auth": map[string]interface{}{
			"client_token": token,
		},
	}
}

func TestAzureLogin_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/azure/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(azureResponse("azure-token-abc"))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	tok, err := AzureLogin(c, "my-role", "jwt-value", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok != "azure-token-abc" {
		t.Errorf("expected azure-token-abc, got %s", tok)
	}
}

func TestAzureLogin_CustomMount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/custom-azure/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(azureResponse("azure-token-xyz"))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	tok, err := AzureLogin(c, "my-role", "jwt-value", "custom-azure")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok != "azure-token-xyz" {
		t.Errorf("expected azure-token-xyz, got %s", tok)
	}
}

func TestAzureLogin_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	_, err := AzureLogin(c, "my-role", "bad-jwt", "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAzureLogin_EmptyToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(azureResponse(""))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	_, err := AzureLogin(c, "my-role", "jwt-value", "")
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}
