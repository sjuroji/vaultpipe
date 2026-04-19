package runner

import (
	"testing"
)

func TestRun_SimpleCommand(t *testing.T) {
	code, err := Run(Options{
		Args:  []string{"echo", "hello"},
		Env:   []string{"PATH=/usr/bin:/bin"},
		Stdin: false,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
}

func TestRun_NonZeroExit(t *testing.T) {
	code, err := Run(Options{
		Args: []string{"sh", "-c", "exit 42"},
		Env:  []string{"PATH=/usr/bin:/bin"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 42 {
		t.Fatalf("expected exit code 42, got %d", code)
	}
}

func TestRun_CommandNotFound(t *testing.T) {
	_, err := Run(Options{
		Args: []string{"__no_such_binary__"},
		Env:  []string{"PATH=/usr/bin:/bin"},
	})
	if err == nil {
		t.Fatal("expected error for missing binary, got nil")
	}
}

func TestRun_NoArgs(t *testing.T) {
	_, err := Run(Options{})
	if err == nil {
		t.Fatal("expected error when no args provided")
	}
}

func TestRun_EnvPassthrough(t *testing.T) {
	code, err := Run(Options{
		Args: []string{"sh", "-c", "test \"$SECRET\" = \"hunter2\""},
		Env:  []string{"PATH=/usr/bin:/bin", "SECRET=hunter2"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if code != 0 {
		t.Fatalf("env var not passed through; exit code %d", code)
	}
}
