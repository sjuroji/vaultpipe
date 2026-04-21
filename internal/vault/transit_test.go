package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func transitEncryptResponse(ciphertext string) map[string]interface{} {
	return map[string]interface{}{
		"data": map[string]string{"ciphertext": ciphertext},
	}
}

func transitDecryptResponse(plaintext string) map[string]interface{} {
	return map[string]interface{}{
		"data": map[string]string{"plaintext": plaintext},
	}
}

func TestEncryptTransit_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/transit/encrypt/mykey" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(transitEncryptResponse("vault:v1:abc123"))
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "tok", HTTP: ts.Client()}
	res, err := EncryptTransit(c, "mykey", "aGVsbG8=")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Ciphertext != "vault:v1:abc123" {
		t.Errorf("expected ciphertext vault:v1:abc123, got %s", res.Ciphertext)
	}
}

func TestEncryptTransit_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "bad", HTTP: ts.Client()}
	_, err := EncryptTransit(c, "mykey", "aGVsbG8=")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDecryptTransit_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/transit/decrypt/mykey" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(transitDecryptResponse("aGVsbG8="))
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "tok", HTTP: ts.Client()}
	res, err := DecryptTransit(c, "mykey", "vault:v1:abc123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Plaintext != "aGVsbG8=" {
		t.Errorf("expected plaintext aGVsbG8=, got %s", res.Plaintext)
	}
}

func TestDecryptTransit_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "tok", HTTP: ts.Client()}
	_, err := DecryptTransit(c, "mykey", "vault:v1:bad")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
