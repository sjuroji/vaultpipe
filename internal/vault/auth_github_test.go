package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func githubResponse(token string, leaseDur int, renewable bool) map[string]interface{} {
	return map[string]interface{}{
		"auth": map[string]interface{}{
			"client_token":   token,
			"lease_duration": leaseDur,
			"renewable":      renewable,
		},
	}
}

func TestGitHubLogin_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/github/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(githubResponse("gh-tok-abc", 3600, true))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	resp, err := GitHubLogin(c, "mygithubpat", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Token != "gh-tok-abc" {
		t.Errorf("expected token gh-tok-abc, got %s", resp.Token)
	}
	if !resp.Renewable {
		t.Error("expected renewable=true")
	}
}

func TestGitHubLogin_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	_, err := GitHubLogin(c, "badtoken", "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGitHubLogin_EmptyToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(githubResponse("", 0, false))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	_, err := GitHubLogin(c, "mygithubpat", "")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestGitHubLogin_CustomMount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/gh-custom/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(githubResponse("custom-tok", 1800, false))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	resp, err := GitHubLogin(c, "mygithubpat", "gh-custom")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Token != "custom-tok" {
		t.Errorf("expected custom-tok, got %s", resp.Token)
	}
}
