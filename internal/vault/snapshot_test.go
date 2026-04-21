package vault

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTakeSnapshot_OK(t *testing.T) {
	payload := []byte("binary-snapshot-data")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Vault-Token") != "test-token" {
			t.Errorf("missing or wrong token")
		}
		w.WriteHeader(http.StatusOK)
		w.Write(payload)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "test-token", HTTP: ts.Client()}
	var buf bytes.Buffer
	n, err := TakeSnapshot(c, &buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != int64(len(payload)) {
		t.Errorf("expected %d bytes, got %d", len(payload), n)
	}
	if !bytes.Equal(buf.Bytes(), payload) {
		t.Errorf("snapshot content mismatch")
	}
}

func TestTakeSnapshot_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "forbidden", http.StatusForbidden)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "bad-token", HTTP: ts.Client()}
	var buf bytes.Buffer
	_, err := TakeSnapshot(c, &buf)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestSnapshotStatus_OK(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	meta := SnapshotMeta{Index: 42, Term: 3, Version: 1, Timestamp: now}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(meta)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "test-token", HTTP: ts.Client()}
	got, err := SnapshotStatus(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Index != 42 {
		t.Errorf("expected index 42, got %d", got.Index)
	}
	if got.Term != 3 {
		t.Errorf("expected term 3, got %d", got.Term)
	}
}

func TestSnapshotStatus_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "test-token", HTTP: ts.Client()}
	_, err := SnapshotStatus(c)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
