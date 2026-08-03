package compose

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yunkhngn/dforge/internal/services"
)

func TestAddServicePreservesComments(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "compose.yaml")
	os.WriteFile(p, []byte("# my stack\nservices:\n  api:\n    build: .\n"), 0o644)

	f, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	svc, ok := services.Get("redis")
	if !ok {
		t.Fatal("redis service not found in catalog")
	}
	if err := f.AddService(svc); err != nil {
		t.Fatal(err)
	}
	out, err := f.Bytes()
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if !strings.Contains(s, "# my stack") {
		t.Fatal("top comment lost")
	}
	if !strings.Contains(s, "redis:") || !strings.Contains(s, "redis:7.4-alpine") {
		t.Fatalf("redis not added:\n%s", s)
	}
}

func TestRemoveServicePrunesVolume(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "compose.yaml")
	f, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	pg, ok := services.Get("postgres")
	if !ok {
		t.Fatal("postgres service not found in catalog")
	}
	if err := f.AddService(pg); err != nil {
		t.Fatal(err)
	}
	if !f.HasService("postgres") {
		t.Fatal("expected HasService(postgres) to be true")
	}
	if err := f.RemoveService("postgres"); err != nil {
		t.Fatal(err)
	}
	if f.HasService("postgres") {
		t.Fatal("expected HasService(postgres) to be false")
	}
	out, err := f.Bytes()
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), "postgres_data") {
		t.Fatalf("volume not pruned:\n%s", out)
	}
}

func TestVolumeRegistration(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "compose.yaml")
	f, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	pg, _ := services.Get("postgres")
	if err := f.AddService(pg); err != nil {
		t.Fatal(err)
	}
	out, err := f.Bytes()
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if !strings.Contains(s, "volumes:") || !strings.Contains(s, "postgres_data:") {
		t.Fatalf("expected top level volume postgres_data registered:\n%s", s)
	}
}

func TestSaveCreatesBackup(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "compose.yaml")
	initialContent := "# initial compose\nservices:\n  app:\n    image: app:latest\n"
	os.WriteFile(p, []byte(initialContent), 0o644)

	f, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	redis, _ := services.Get("redis")
	if err := f.AddService(redis); err != nil {
		t.Fatal(err)
	}
	if err := f.Save(); err != nil {
		t.Fatal(err)
	}

	// Verify backup file exists and matches initial content
	bakPath := p + ".bak"
	bakBytes, err := os.ReadFile(bakPath)
	if err != nil {
		t.Fatalf("expected backup file %s to exist: %v", bakPath, err)
	}
	if string(bakBytes) != initialContent {
		t.Fatalf("backup content mismatch.\nExpected:\n%s\nGot:\n%s", initialContent, string(bakBytes))
	}

	// Verify current file has new content
	newBytes, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(newBytes), "redis:") {
		t.Fatalf("new compose file missing redis:\n%s", string(newBytes))
	}
}

func TestLoadNonExistentFile(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "new-compose.yaml")
	f, err := Load(p)
	if err != nil {
		t.Fatalf("Load non-existent file failed: %v", err)
	}
	if f.HasService("anything") {
		t.Fatal("expected HasService to be false for empty file")
	}
	redis, _ := services.Get("redis")
	if err := f.AddService(redis); err != nil {
		t.Fatal(err)
	}
	out, err := f.Bytes()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), "redis:") {
		t.Fatalf("expected output to contain redis:\n%s", string(out))
	}
}

func TestErrorsAndEdgeCases(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "compose.yaml")
	f, _ := Load(p)

	// Remove non-existent service
	if err := f.RemoveService("nonexistent"); err == nil {
		t.Fatal("expected error removing non-existent service")
	}

	// Add service twice
	redis, _ := services.Get("redis")
	if err := f.AddService(redis); err != nil {
		t.Fatal(err)
	}
	if err := f.AddService(redis); err == nil {
		t.Fatal("expected error adding existing service")
	}
}
