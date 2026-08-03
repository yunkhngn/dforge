package env

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCompareReportsDiffs(t *testing.T) {
	dir := t.TempDir()
	env := filepath.Join(dir, ".env")
	ex := filepath.Join(dir, ".env.example")
	os.WriteFile(ex, []byte("A=1\nB=2\nC=3\n"), 0o644)
	os.WriteFile(env, []byte("A=1\nA=9\nB=\nD=4\n"), 0o644)

	r, err := Compare(env, ex)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Missing) != 1 || r.Missing[0] != "C" {
		t.Fatalf("missing: %v, expected [C]", r.Missing)
	}
	if len(r.Extra) != 1 || r.Extra[0] != "D" {
		t.Fatalf("extra: %v, expected [D]", r.Extra)
	}
	if len(r.Duplicates) != 1 || r.Duplicates[0] != "A" {
		t.Fatalf("dups: %v, expected [A]", r.Duplicates)
	}
	if len(r.Empty) != 1 || r.Empty[0] != "B" {
		t.Fatalf("empty: %v, expected [B]", r.Empty)
	}
}

func TestCompareMissingEnvFile(t *testing.T) {
	dir := t.TempDir()
	ex := filepath.Join(dir, ".env.example")
	os.WriteFile(ex, []byte("A=1\n"), 0o644)
	r, err := Compare(filepath.Join(dir, ".env"), ex)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.MissingFile {
		t.Fatal("expected MissingFile=true")
	}
}

func TestCompareMultipleDuplicatesAndComments(t *testing.T) {
	dir := t.TempDir()
	env := filepath.Join(dir, ".env")
	ex := filepath.Join(dir, ".env.example")
	os.WriteFile(ex, []byte("# Example env file\nA=1\nB=2\n"), 0o644)
	os.WriteFile(env, []byte("# Comments\n\nA=1\nA=2\nA=3\nB=2\n"), 0o644)

	r, err := Compare(env, ex)
	if err != nil {
		t.Fatal(err)
	}
	if r.MissingFile {
		t.Fatal("expected MissingFile=false")
	}
	if len(r.Duplicates) != 1 || r.Duplicates[0] != "A" {
		t.Fatalf("dups: %v, expected [A]", r.Duplicates)
	}
	if len(r.Missing) != 0 {
		t.Fatalf("missing: %v, expected []", r.Missing)
	}
	if len(r.Extra) != 0 {
		t.Fatalf("extra: %v, expected []", r.Extra)
	}
	if len(r.Empty) != 0 {
		t.Fatalf("empty: %v, expected []", r.Empty)
	}
}

func TestCompareMissingExampleFile(t *testing.T) {
	dir := t.TempDir()
	env := filepath.Join(dir, ".env")
	os.WriteFile(env, []byte("A=1\n"), 0o644)
	_, err := Compare(env, filepath.Join(dir, ".env.example"))
	if err == nil {
		t.Fatal("expected error when example file is missing")
	}
}

func TestCompareNoDiffs(t *testing.T) {
	dir := t.TempDir()
	env := filepath.Join(dir, ".env")
	ex := filepath.Join(dir, ".env.example")
	os.WriteFile(ex, []byte("A=1\nB=2\n"), 0o644)
	os.WriteFile(env, []byte("A=10\nB=20\n"), 0o644)

	r, err := Compare(env, ex)
	if err != nil {
		t.Fatal(err)
	}
	if r.MissingFile {
		t.Fatal("expected MissingFile=false")
	}
	if len(r.Missing) != 0 || len(r.Extra) != 0 || len(r.Duplicates) != 0 || len(r.Empty) != 0 {
		t.Fatalf("expected no diffs, got: %+v", r)
	}
}
