package envfile

import (
	"os"
	"testing"
)

func writeTempEnvFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.env")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestParse_BasicKeyValue(t *testing.T) {
	path := writeTempEnvFile(t, "FOO=bar\nBAZ=qux\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["FOO"] != "bar" || env["BAZ"] != "qux" {
		t.Errorf("unexpected env: %v", env)
	}
}

func TestParse_QuotedValues(t *testing.T) {
	path := writeTempEnvFile(t, `SINGLE='hello world'
DOUBLE="goodbye world"
`)
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["SINGLE"] != "hello world" {
		t.Errorf("SINGLE: got %q", env["SINGLE"])
	}
	if env["DOUBLE"] != "goodbye world" {
		t.Errorf("DOUBLE: got %q", env["DOUBLE"])
	}
}

func TestParse_CommentsAndBlanks(t *testing.T) {
	path := writeTempEnvFile(t, "# this is a comment\n\nKEY=value\n")
	env, err := Parse(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(env) != 1 || env["KEY"] != "value" {
		t.Errorf("unexpected env: %v", env)
	}
}

func TestParse_InvalidLine(t *testing.T) {
	path := writeTempEnvFile(t, "NOTVALID\n")
	_, err := Parse(path)
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestParse_MissingFile(t *testing.T) {
	_, err := Parse("/nonexistent/path/.env")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
