package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func rekeyStatusResponse(started bool, nonce string) map[string]interface{} {
	return map[string]interface{}{
		"started":               started,
		"nonce":                 nonce,
		"t":                     3,
		"n":                     5,
		"progress":              0,
		"required":              3,
		"pgp_fingerprints":      []string{},
		"backup":                false,
		"verification_required": false,
	}
}

func TestInitRekey_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/sys/rekey/init" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(rekeyStatusResponse(true, "abc-nonce"))
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "root", HTTP: ts.Client()}
	status, err := InitRekey(c, 5, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.Started {
		t.Error("expected started=true")
	}
	if status.Nonce != "abc-nonce" {
		t.Errorf("expected nonce abc-nonce, got %s", status.Nonce)
	}
}

func TestInitRekey_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "bad", HTTP: ts.Client()}
	_, err := InitRekey(c, 5, 3)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGetRekeyStatus_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(rekeyStatusResponse(false, ""))
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "root", HTTP: ts.Client()}
	status, err := GetRekeyStatus(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Started {
		t.Error("expected started=false")
	}
}

func TestCancelRekey_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "root", HTTP: ts.Client()}
	if err := CancelRekey(c); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCancelRekey_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "root", HTTP: ts.Client()}
	if err := CancelRekey(c); err == nil {
		t.Fatal("expected error, got nil")
	}
}
