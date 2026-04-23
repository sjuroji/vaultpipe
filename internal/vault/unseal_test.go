package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func unsealResponse(sealed bool, progress, threshold int) map[string]interface{} {
	return map[string]interface{}{
		"sealed":    sealed,
		"t":         threshold,
		"n":         5,
		"progress":  progress,
		"threshold": threshold,
	}
}

func TestSubmitUnsealKey_OK(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/v1/sys/unseal" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(unsealResponse(true, 1, 3))
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "root", HTTP: ts.Client()}
	res, err := SubmitUnsealKey(c, "abc123", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Sealed {
		t.Error("expected sealed=true")
	}
	if res.Progress != 1 {
		t.Errorf("expected progress=1, got %d", res.Progress)
	}
}

func TestSubmitUnsealKey_Reset(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req UnsealRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if !req.Reset {
			t.Error("expected reset=true in request")
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(unsealResponse(true, 0, 3))
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "root", HTTP: ts.Client()}
	res, err := SubmitUnsealKey(c, "", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Progress != 0 {
		t.Errorf("expected progress=0 after reset, got %d", res.Progress)
	}
}

func TestSubmitUnsealKey_NonOKStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	c := &Client{Address: ts.URL, Token: "root", HTTP: ts.Client()}
	_, err := SubmitUnsealKey(c, "badkey", false)
	if err == nil {
		t.Fatal("expected error for non-OK status")
	}
}
