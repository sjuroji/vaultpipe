package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func kubernetesResponse(token string) map[string]interface{} {
	return map[string]interface{}{
		"auth": map[string]interface{}{
			"client_token": token,
		},
	}
}

func TestKubernetesLogin_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/kubernetes/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(kubernetesResponse("k8s-token-abc"))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	tok, err := KubernetesLogin(c, "my-role", "jwt-value", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok != "k8s-token-abc" {
		t.Errorf("expected k8s-token-abc, got %s", tok)
	}
}

func TestKubernetesLogin_CustomMount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/k8s-prod/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(kubernetesResponse("prod-token"))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	tok, err := KubernetesLogin(c, "prod-role", "jwt-value", "k8s-prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok != "prod-token" {
		t.Errorf("expected prod-token, got %s", tok)
	}
}

func TestKubernetesLogin_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	_, err := KubernetesLogin(c, "role", "jwt", "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestKubernetesLogin_EmptyToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(kubernetesResponse(""))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	_, err := KubernetesLogin(c, "role", "jwt", "")
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}
