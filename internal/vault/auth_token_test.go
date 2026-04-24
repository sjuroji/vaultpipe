package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func tokenAuthResponse(id, accessor string, policies []string, renewable bool, ttl int) map[string]interface{} {
	return map[string]interface{}{
		"data": map[string]interface{}{
			"id":        id,
			"accessor":  accessor,
			"policies":  policies,
			"renewable": renewable,
			"ttl":       ttl,
		},
	}
}

func TestTokenLogin_OK(t *testing.T) {
	payload := tokenAuthResponse("mytoken", "acc123", []string{"default", "dev"}, true, 3600)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != "mytoken" {
			t.Errorf("expected token header, got %q", r.Header.Get("X-Vault-Token"))
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(payload)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, HTTP: ts.Client()}
	res, err := TokenLogin(c, "mytoken")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ClientToken != "mytoken" {
		t.Errorf("expected ClientToken mytoken, got %q", res.ClientToken)
	}
	if res.Accessor != "acc123" {
		t.Errorf("expected accessor acc123, got %q", res.Accessor)
	}
	if len(res.Policies) != 2 {
		t.Errorf("expected 2 policies, got %d", len(res.Policies))
	}
	if !res.Renewable {
		t.Error("expected renewable to be true")
	}
	if res.TTL != 3600 {
		t.Errorf("expected TTL 3600, got %d", res.TTL)
	}
}

func TestTokenLogin_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, HTTP: ts.Client()}
	_, err := TokenLogin(c, "badtoken")
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}
