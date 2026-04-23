package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func tokenLookupResponse(id, display string, policies []string, ttl int, renewable bool, creationTime int64) map[string]interface{} {
	return map[string]interface{}{
		"data": map[string]interface{}{
			"id":            id,
			"display_name":  display,
			"policies":      policies,
			"ttl":           ttl,
			"renewable":     renewable,
			"creation_time": creationTime,
		},
	}
}

func TestLookupSelfToken_OK(t *testing.T) {
	policies := []string{"default", "read-secrets"}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/token/lookup-self" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tokenLookupResponse("tok123", "my-token", policies, 3600, true, 1700000000))
	}))
	defer ts.Close()

	client := &Client{Address: ts.URL, Token: "tok123", HTTP: ts.Client()}
	info, err := LookupSelfToken(client)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.ID != "tok123" {
		t.Errorf("expected ID tok123, got %s", info.ID)
	}
	if info.DisplayName != "my-token" {
		t.Errorf("expected display name my-token, got %s", info.DisplayName)
	}
	if len(info.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(info.Policies))
	}
	if info.TTL != 3600*time.Second {
		t.Errorf("expected TTL 3600s, got %v", info.TTL)
	}
	if !info.Renewable {
		t.Error("expected renewable to be true")
	}
	if info.CreationTime != time.Unix(1700000000, 0) {
		t.Errorf("unexpected creation time: %v", info.CreationTime)
	}
}

func TestLookupSelfToken_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	client := &Client{Address: ts.URL, Token: "bad", HTTP: ts.Client()}
	_, err := LookupSelfToken(client)
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}

func TestLookupSelfToken_InvalidJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{invalid json`))
	}))
	defer ts.Close()

	client := &Client{Address: ts.URL, Token: "tok", HTTP: ts.Client()}
	_, err := LookupSelfToken(client)
	if err == nil {
		t.Fatal("expected error for invalid JSON response")
	}
}

func TestTokenInfo_IsExpired(t *testing.T) {
	expired := &TokenInfo{TTL: 0}
	if !expired.IsExpired() {
		t.Error("expected IsExpired true for TTL=0")
	}

	active := &TokenInfo{TTL: 60 * time.Second}
	if active.IsExpired() {
		t.Error("expected IsExpired false for TTL=60s")
	}
}
