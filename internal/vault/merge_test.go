package vault

import (
	"testing"
)

func TestMergeIntoEnv_OverridesExisting(t *testing.T) {
	base := map[string]string{
		"FOO": "original",
		"BAR": "keep",
	}
	secrets := map[string]string{
		"FOO": "overridden",
		"NEW": "added",
	}
	result := MergeIntoEnv(base, secrets)
	if result["FOO"] != "overridden" {
		t.Errorf("expected FOO=overridden, got %s", result["FOO"])
	}
	if result["BAR"] != "keep" {
		t.Errorf("expected BAR=keep, got %s", result["BAR"])
	}
	if result["NEW"] != "added" {
		t.Errorf("expected NEW=added, got %s", result["NEW"])
	}
}

func TestMergeIntoEnv_EmptySecrets(t *testing.T) {
	base := map[string]string{"A": "1"}
	result := MergeIntoEnv(base, map[string]string{})
	if result["A"] != "1" {
		t.Errorf("expected A=1, got %s", result["A"])
	}
}

func TestToSlice_Format(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	slice := ToSlice(env)
	if len(slice) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(slice))
	}
	if slice[0] != "KEY=value" {
		t.Errorf("unexpected format: %s", slice[0])
	}
}

func TestFromSlice_ParsesCorrectly(t *testing.T) {
	slice := []string{"FOO=bar", "BAZ=qux"}
	result := FromSlice(slice)
	if result["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %s", result["FOO"])
	}
	if result["BAZ"] != "qux" {
		t.Errorf("expected BAZ=qux, got %s", result["BAZ"])
	}
}

func TestFromSlice_SkipsMalformed(t *testing.T) {
	slice := []string{"VALID=yes", "NOEQUALSSIGN", ""}
	result := FromSlice(slice)
	if _, ok := result["NOEQUALSSIGN"]; ok {
		t.Error("expected malformed entry to be skipped")
	}
	if result["VALID"] != "yes" {
		t.Errorf("expected VALID=yes, got %s", result["VALID"])
	}
}
