package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func ociResponse(token string) map[string]interface{} {
	return map[string]interface{}{
		"auth": map[string]interface{}{
			"client_token": token,
		},
	}
}

func TestOCILogin_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/oci/login/my-role" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ociResponse("oci-token-abc"))
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	headers := map[string]string{"date": "Mon, 01 Jan 2024 00:00:00 GMT"}
	resp, err := OCILogin(c, "my-role", headers, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Token != "oci-token-abc" {
		t.Errorf("expected token 'oci-token-abc', got %q", resp.Token)
	}
}

func TestOCILogin_CustomMount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/custom-oci/login/my-role" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ociResponse("oci-custom-token"))
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	resp, err := OCILogin(c, "my-role", nil, "custom-oci")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Token != "oci-custom-token" {
		t.Errorf("expected 'oci-custom-token', got %q", resp.Token)
	}
}

func TestOCILogin_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	_, err := OCILogin(c, "my-role", nil, "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestOCILogin_EmptyToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ociResponse(""))
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	_, err := OCILogin(c, "my-role", nil, "")
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}
