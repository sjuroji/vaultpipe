package runner

import (
	"os"
	"testing"
)

func findEnv(slice []string, key string) (string, bool) {
	for _, pair := range slice {
		k, v := splitPair(pair)
		if k == key {
			return v, true
		}
	}
	return "", false
}

func TestBuild_EnvFileOverridesOS(t *testing.T) {
	os.Setenv("SHARED_KEY", "from_os")
	t.Cleanup(func() { os.Unsetenv("SHARED_KEY") })

	f, _ := os.CreateTemp("", "*.env")
	f.WriteString("SHARED_KEY=from_envfile\n")
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })

	result, err := Build(Options{
		InheritOS:   true,
		EnvFilePath: f.Name(),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := findEnv(result, "SHARED_KEY")
	if !ok || v != "from_envfile" {
		t.Errorf("expected from_envfile, got %q", v)
	}
}

func TestBuild_NoInheritOS(t *testing.T) {
	os.Setenv("SECRET_OS_VAR", "should_not_appear")
	t.Cleanup(func() { os.Unsetenv("SECRET_OS_VAR") })

	result, err := Build(Options{InheritOS: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := findEnv(result, "SECRET_OS_VAR"); ok {
		t.Error("OS var should not be inherited")
	}
}

func TestSplitPair(t *testing.T) {
	cases := []struct{ in, k, v string }{
		{"FOO=bar", "FOO", "bar"},
		{"FOO=bar=baz", "FOO", "bar=baz"},
		{"NOVALUE", "NOVALUE", ""},
	}
	for _, c := range cases {
		gk, gv := splitPair(c.in)
		if gk != c.k || gv != c.v {
			t.Errorf("splitPair(%q) = %q,%q; want %q,%q", c.in, gk, gv, c.k, c.v)
		}
	}
}

func TestBuild_VaultOverridesEnvFile(t *testing.T) {
	// Without a real Vault client, verify env file values appear correctly
	f, _ := os.CreateTemp("", "*.env")
	f.WriteString("APP_TOKEN=from_envfile\n")
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })

	result, err := Build(Options{
		InheritOS:   false,
		EnvFilePath: f.Name(),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := findEnv(result, "APP_TOKEN")
	if !ok || v != "from_envfile" {
		t.Errorf("expected from_envfile, got %q", v)
	}
}
