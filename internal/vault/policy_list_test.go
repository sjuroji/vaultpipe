package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func policyListResponse(keys []string) map[string]interface{} {
	return map[string]interface{}{
		"data": map[string]interface{}{
			"keys": keys,
		},
	}
}

func TestListPolicies_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(policyListResponse([]string{"default", "admin", "readonly"}))
	}))
	defer ts.Close()

	client, _ := NewClient(ts.URL, "test-token")
	result, err := ListPolicies(client)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result.Policies) != 3 {
		t.Errorf("expected 3 policies, got %d", len(result.Policies))
	}
	if result.Policies[1] != "admin" {
		t.Errorf("expected admin, got %s", result.Policies[1])
	}
}

func TestListPolicies_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	client, _ := NewClient(ts.URL, "test-token")
	_, err := ListPolicies(client)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetPolicy_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": map[string]interface{}{
				"name":  "admin",
				"rules": `path "secret/*" { capabilities = ["read"] }`,
			},
		})
	}))
	defer ts.Close()

	client, _ := NewClient(ts.URL, "test-token")
	policy, err := GetPolicy(client, "admin")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if policy.Name != "admin" {
		t.Errorf("expected admin, got %s", policy.Name)
	}
	if policy.Rules == "" {
		t.Error("expected non-empty rules")
	}
}

func TestGetPolicy_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	client, _ := NewClient(ts.URL, "test-token")
	_, err := GetPolicy(client, "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
