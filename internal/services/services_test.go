package services

import (
	"strings"
	"testing"
)

func TestGetPostgresIsPinned(t *testing.T) {
	svc, ok := Get("postgres")
	if !ok {
		t.Fatal("postgres not in catalog")
	}
	if !strings.Contains(svc.Image, ":") || strings.HasSuffix(svc.Image, ":latest") {
		t.Fatalf("postgres image must be pinned, got %q", svc.Image)
	}
	if svc.Healthcheck == nil {
		t.Fatal("postgres must define a healthcheck")
	}
}

func TestNamesCoversBrief(t *testing.T) {
	want := []string{"mailpit", "meilisearch", "minio", "mongodb", "mysql", "postgres", "rabbitmq", "redis"}
	got := Names()
	if strings.Join(got, ",") != strings.Join(want, ",") {
		t.Fatalf("want %v, got %v", want, got)
	}
}

func TestGetUnknownService(t *testing.T) {
	_, ok := Get("nonexistent")
	if ok {
		t.Fatal("expected false for unknown service, got true")
	}
}

func TestAllCatalogImagesArePinned(t *testing.T) {
	names := Names()
	if len(names) == 0 {
		t.Fatal("catalog is empty")
	}
	for _, name := range names {
		svc, ok := Get(name)
		if !ok {
			t.Fatalf("service %s listed in Names() but not returned by Get()", name)
		}
		if svc.Name != name {
			t.Errorf("service name mismatch: got %q, want %q", svc.Name, name)
		}
		if svc.Image == "" {
			t.Errorf("service %s has empty image", name)
		}
		if !strings.Contains(svc.Image, ":") || strings.HasSuffix(svc.Image, ":latest") {
			t.Errorf("service %s image must be pinned explicitly, got %q", name, svc.Image)
		}
	}
}
