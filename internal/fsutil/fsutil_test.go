package fsutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExists(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "exists.txt")

	if Exists(p) {
		t.Fatalf("expected Exists(%q) to be false", p)
	}

	if err := os.WriteFile(p, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	if !Exists(p) {
		t.Fatalf("expected Exists(%q) to be true", p)
	}
}

func TestWriteNew(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "f.txt")

	// First write should succeed
	if err := WriteNew(p, []byte("content")); err != nil {
		t.Fatalf("expected WriteNew to succeed, got: %v", err)
	}

	cur, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if string(cur) != "content" {
		t.Fatalf("content wrong: %q", string(cur))
	}

	// Second write should fail because file exists
	if err := WriteNew(p, []byte("new")); err == nil {
		t.Fatal("expected error writing over existing file")
	}
}

func TestWriteWithBackup(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "compose.yaml")

	// First write when file does not exist
	if err := WriteWithBackup(p, []byte("initial")); err != nil {
		t.Fatalf("expected WriteWithBackup on new file to succeed, got: %v", err)
	}

	cur, _ := os.ReadFile(p)
	if string(cur) != "initial" {
		t.Fatalf("initial content wrong: %q", string(cur))
	}

	// Second write when file exists -> should create backup and overwrite
	if err := WriteWithBackup(p, []byte("updated")); err != nil {
		t.Fatalf("expected WriteWithBackup on existing file to succeed, got: %v", err)
	}

	bak, err := os.ReadFile(p + ".bak")
	if err != nil {
		t.Fatalf("failed to read backup file: %v", err)
	}
	if string(bak) != "initial" {
		t.Fatalf("backup wrong: %q", string(bak))
	}

	cur, _ = os.ReadFile(p)
	if string(cur) != "updated" {
		t.Fatalf("current wrong: %q", string(cur))
	}
}
