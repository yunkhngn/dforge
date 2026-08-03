package detect

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestDetectNextjs(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "package.json", `{"dependencies":{"next":"14.0.0","react":"18.0.0"}}`)
	got, err := Detect(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Framework != Nextjs {
		t.Fatalf("want [nextjs], got %+v", got)
	}
}

func TestDetectGo(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "go.mod", "module example.com/x\n\ngo 1.22\n")
	got, _ := Detect(dir)
	if len(got) != 1 || got[0].Framework != Go {
		t.Fatalf("want [go], got %+v", got)
	}
}

func TestDetectNoneReturnsEmpty(t *testing.T) {
	got, err := Detect(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Fatalf("want empty, got %+v", got)
	}
}

func TestDetectNestjs(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "package.json", `{"dependencies":{"@nestjs/core":"10.0.0"}}`)
	got, err := Detect(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Framework != Nestjs {
		t.Fatalf("want [nestjs], got %+v", got)
	}
}

func TestDetectExpress(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "package.json", `{"dependencies":{"express":"4.18.0"}}`)
	got, err := Detect(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Framework != Express {
		t.Fatalf("want [express], got %+v", got)
	}
}

func TestDetectReact(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "package.json", `{"dependencies":{"react":"18.0.0"}}`)
	got, err := Detect(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Framework != React {
		t.Fatalf("want [react], got %+v", got)
	}
}

func TestDetectNodeDefault(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "package.json", `{"dependencies":{"lodash":"4.17.21"}}`)
	got, err := Detect(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Framework != Express {
		t.Fatalf("want [express] default for node, got %+v", got)
	}
}

func TestDetectRust(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "Cargo.toml", "[package]\nname = \"foo\"\n")
	got, err := Detect(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0].Framework != Rust {
		t.Fatalf("want [rust], got %+v", got)
	}
}

func TestDetectSpringBoot(t *testing.T) {
	tests := []string{"pom.xml", "build.gradle", "build.gradle.kts"}
	for _, file := range tests {
		t.Run(file, func(t *testing.T) {
			dir := t.TempDir()
			writeFile(t, dir, file, "")
			got, err := Detect(dir)
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != 1 || got[0].Framework != SpringBoot {
				t.Fatalf("want [springboot], got %+v", got)
			}
		})
	}
}

func TestDetectInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	writeFile(t, dir, "package.json", `{invalid json`)
	_, err := Detect(dir)
	if err == nil {
		t.Fatal("expected error for invalid json, got nil")
	}
}
