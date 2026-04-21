package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func raftConfigResponse(index int, servers []RaftServer) map[string]interface{} {
	return map[string]interface{}{
		"data": map[string]interface{}{
			"index":   index,
			"servers": servers,
		},
	}
}

func TestGetRaftConfiguration_OK(t *testing.T) {
	servers := []RaftServer{
		{ID: "node1", Address: "127.0.0.1:8201", Leader: true, Voter: true, Protocol: "3"},
		{ID: "node2", Address: "127.0.0.1:8202", Leader: false, Voter: true, Protocol: "3"},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/storage/raft/configuration" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(raftConfigResponse(1, servers))
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "test-token", HTTP: ts.Client()}
	cfg, err := GetRaftConfiguration(c)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Servers) != 2 {
		t.Errorf("expected 2 servers, got %d", len(cfg.Servers))
	}
	if !cfg.Servers[0].Leader {
		t.Errorf("expected first server to be leader")
	}
	if cfg.Servers[0].ID != "node1" {
		t.Errorf("expected node1, got %s", cfg.Servers[0].ID)
	}
}

func TestGetRaftConfiguration_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "bad-token", HTTP: ts.Client()}
	_, err := GetRaftConfiguration(c)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRemoveRaftPeer_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/sys/storage/raft/remove-peer" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["server_id"] != "node2" {
			t.Errorf("expected node2, got %s", body["server_id"])
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "test-token", HTTP: ts.Client()}
	if err := RemoveRaftPeer(c, "node2"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRemoveRaftPeer_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "test-token", HTTP: ts.Client()}
	if err := RemoveRaftPeer(c, "node-missing"); err == nil {
		t.Fatal("expected error, got nil")
	}
}
