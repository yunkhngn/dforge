package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yunkhngn/dforge/internal/detect"
)

func TestRunInitWritesArtifacts(t *testing.T) {
	dir := t.TempDir()
	if err := runInit(dir, detect.Go, true); err != nil {
		t.Fatal(err)
	}
	for _, f := range []string{"Dockerfile", "compose.yaml", ".dockerignore"} {
		if _, err := os.Stat(filepath.Join(dir, f)); err != nil {
			t.Fatalf("missing generated %s", f)
		}
	}
}

func TestRunInitRefusesOverwriteWithoutForce(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "Dockerfile"), []byte("keep"), 0o644)
	err := runInit(dir, detect.Go, false)
	if err == nil {
		t.Fatal("expected refusal to overwrite existing Dockerfile")
	}
	// The existing file must be untouched (no .bak from a partial write)...
	if content, _ := os.ReadFile(filepath.Join(dir, "Dockerfile")); string(content) != "keep" {
		t.Fatalf("existing Dockerfile was modified: %q", content)
	}
	// ...and no other artifact should have been written.
	for _, f := range []string{"compose.yaml", ".dockerignore", "Dockerfile.bak"} {
		if _, err := os.Stat(filepath.Join(dir, f)); err == nil {
			t.Fatalf("partial write occurred: %s should not exist", f)
		}
	}
}

func TestPickFramework(t *testing.T) {
	t.Run("override flag takes precedence", func(t *testing.T) {
		fw, err := pickFramework([]detect.Candidate{{Framework: detect.Go}}, detect.React)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fw != detect.React {
			t.Fatalf("expected react, got %s", fw)
		}
	})

	t.Run("single candidate auto-selects", func(t *testing.T) {
		fw, err := pickFramework([]detect.Candidate{{Framework: detect.Go}}, "")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if fw != detect.Go {
			t.Fatalf("expected go, got %s", fw)
		}
	})

	t.Run("no candidate returns error", func(t *testing.T) {
		_, err := pickFramework(nil, "")
		if err == nil {
			t.Fatal("expected error when no candidates found")
		}
	})
}
