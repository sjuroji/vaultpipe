package runner

import (
	"strings"
	"testing"
)

func findEnv(env []string, key string) (string, bool) {
	prefix := key + "="
	for _, e := range env {
		if strings.HasPrefix(e, prefix) {
			return strings.TrimPrefix(e, prefix), true
		}
	}
	return "", false
}

func TestBuild_VaultOverridesEnvFile(t *testing.T) {
	src := EnvSource{
		VaultSecrets: map[string]string{"DB_PASS": "vault-secret"},
		EnvFileVars:  map[string]string{"DB_PASS": "file-secret"},
		InheritOS:    false,
	}
	env := Build(src)
	v, ok := findEnv(env, "DB_PASS")
	if !ok {
		t.Fatal("DB_PASS not found in env")
	}
	if v != "vault-secret" {
		t.Fatalf("expected vault-secret, got %s", v)
	}
}

func TestBuild_EnvFileOverridesOS(t *testing.T) {
	t.Setenv("MY_VAR", "os-value")
	src := EnvSource{
		VaultSecrets: map[string]string{},
		EnvFileVars:  map[string]string{"MY_VAR": "file-value"},
		InheritOS:    true,
	}
	env := Build(src)
	v, ok := findEnv(env, "MY_VAR")
	if !ok {
		t.Fatal("MY_VAR not found")
	}
	if v != "file-value" {
		t.Fatalf("expected file-value, got %s", v)
	}
}

func TestBuild_NoInheritOS(t *testing.T) {
	t.Setenv("SHOULD_NOT_APPEAR", "yes")
	src := EnvSource{
		InheritOS: false,
	}
	env := Build(src)
	if _, ok := findEnv(env, "SHOULD_NOT_APPEAR"); ok {
		t.Fatal("OS env should not be inherited")
	}
}

func TestSplitPair(t *testing.T) {
	k, v := splitPair("FOO=bar=baz")
	if k != "FOO" || v != "bar=baz" {
		t.Fatalf("unexpected split: %q %q", k, v)
	}
}
