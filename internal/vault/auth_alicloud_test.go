package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func alicloudResponse(token string) map[string]interface{} {
	return map[string]interface{}{
		"auth": map[string]interface{}{
			"client_token": token,
		},
	}
}

func TestAliCloudLogin_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/alicloud/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(alicloudResponse("tok-alicloud-abc"))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	resp, err := AliCloudLogin(c, "my-role", "https://sts.example.com", "{}", "")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Token != "tok-alicloud-abc" {
		t.Errorf("expected token tok-alicloud-abc, got %s", resp.Token)
	}
}

func TestAliCloudLogin_CustomMount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/custom-ali/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(alicloudResponse("tok-custom"))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	resp, err := AliCloudLogin(c, "role", "https://sts.example.com", "{}", "custom-ali")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Token != "tok-custom" {
		t.Errorf("expected tok-custom, got %s", resp.Token)
	}
}

func TestAliCloudLogin_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	_, err := AliCloudLogin(c, "role", "https://sts.example.com", "{}", "")
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestAliCloudLogin_EmptyToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(alicloudResponse(""))
	}))
	defer ts.Close()

	c := newTestClient(ts.URL)
	_, err := AliCloudLogin(c, "role", "https://sts.example.com", "{}", "")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}
