package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func awsResponse(token string) map[string]interface{} {
	return map[string]interface{}{
		"auth": map[string]interface{}{
			"client_token": token,
		},
	}
}

func TestAWSLogin_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/aws/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(awsResponse("aws-token-abc"))
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	token, err := AWSLogin(c, "my-role", "aHR0cHM6Ly9zdHM=", "QWN0aW9uPUdldENhbGxlcklkZW50aXR5", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "aws-token-abc" {
		t.Errorf("expected aws-token-abc, got %s", token)
	}
}

func TestAWSLogin_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	_, err := AWSLogin(c, "role", "url", "body", "headers", "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAWSLogin_EmptyToken(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(awsResponse(""))
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	_, err := AWSLogin(c, "role", "url", "body", "headers", "")
	if err == nil {
		t.Fatal("expected error for empty token")
	}
}

func TestAWSLogin_CustomMount(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/auth/custom-aws/login" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(awsResponse("tok"))
	}))
	defer ts.Close()

	c := newTestClient(t, ts.URL)
	token, err := AWSLogin(c, "role", "url", "body", "headers", "custom-aws")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "tok" {
		t.Errorf("expected tok, got %s", token)
	}
}
