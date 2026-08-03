package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunAddThenRemove(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "compose.yaml")
	os.WriteFile(p, []byte("services:\n  api:\n    build: .\n"), 0o644)

	if err := runAdd(p, "redis"); err != nil {
		t.Fatal(err)
	}
	out, _ := os.ReadFile(p)
	if !strings.Contains(string(out), "redis:") {
		t.Fatalf("redis not added:\n%s", out)
	}
	if err := runRemove(p, "redis"); err != nil {
		t.Fatal(err)
	}
	out, _ = os.ReadFile(p)
	if strings.Contains(string(out), "redis:7.4-alpine") {
		t.Fatalf("redis not removed:\n%s", out)
	}
}

func TestRunAddUnknownService(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "compose.yaml")
	os.WriteFile(p, []byte("services: {}\n"), 0o644)
	if err := runAdd(p, "cassandra"); err == nil {
		t.Fatal("expected error for unknown service")
	}
}
