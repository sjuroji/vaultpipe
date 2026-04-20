package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func auditBackendsResponse() map[string]interface{} {
	return map[string]interface{}{
		"data": map[string]interface{}{
			"file/": map[string]interface{}{
				"type":        "file",
				"description": "file audit backend",
			},
			"syslog/": map[string]interface{}{
				"type":        "syslog",
				"description": "syslog audit backend",
			},
		},
	}
}

func TestListAuditBackends_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/audit" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(auditBackendsResponse())
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "test-token", HTTP: ts.Client()}
	backends, err := ListAuditBackends(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(backends) != 2 {
		t.Errorf("expected 2 backends, got %d", len(backends))
	}
	if backends["file/"].Type != "file" {
		t.Errorf("expected file backend type, got %s", backends["file/"].Type)
	}
}

func TestListAuditBackends_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "bad-token", HTTP: ts.Client()}
	_, err := ListAuditBackends(c)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAuditLog_Record(t *testing.T) {
	log := &AuditLog{}
	before := time.Now().UTC()
	log.Record("secret/data/myapp", "read", true)
	log.Record("secret/data/missing", "read", false)

	if len(log.Events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(log.Events))
	}
	if !log.Events[0].Success {
		t.Error("expected first event to be successful")
	}
	if log.Events[1].Success {
		t.Error("expected second event to be unsuccessful")
	}
	if log.Events[0].Timestamp.Before(before) {
		t.Error("timestamp should be after test start")
	}
}
