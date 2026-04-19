package cmd

import (
	"bytes"
	"os"
	"testing"
)

func TestExecute_NoArgs_PrintsUsage(t *testing.T) {
	// Capture stderr by redirecting rootCmd output
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{})

	err := rootCmd.Execute()
	if err == nil {
		t.Fatal("expected error for no args, got nil")
	}
}

func TestExecute_HelpFlag(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--help"})

	// --help exits cleanly
	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Fatal("expected help output, got nothing")
	}
}

func TestFlags_Defaults(t *testing.T) {
	if envFile != "" {
		t.Errorf("expected empty envFile default, got %q", envFile)
	}
	if vaultPath != "" {
		t.Errorf("expected empty vaultPath default, got %q", vaultPath)
	}
	if !inheritEnv {
		t.Error("expected inheritEnv to default to true")
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
