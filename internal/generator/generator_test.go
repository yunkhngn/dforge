package generator

import (
	"strings"
	"testing"

	"github.com/yunkhngn/dforge/internal/detect"
)

func TestRenderAllFrameworksProduceThreeFiles(t *testing.T) {
	frameworks := []detect.Framework{
		detect.Nextjs,
		detect.React,
		detect.Nestjs,
		detect.Express,
		detect.SpringBoot,
		detect.Go,
		detect.Rust,
	}

	for _, fw := range frameworks {
		t.Run(string(fw), func(t *testing.T) {
			out, err := Render(fw)
			if err != nil {
				t.Fatalf("unexpected error rendering %s: %v", fw, err)
			}

			for _, key := range []string{"Dockerfile", "compose.yaml", ".dockerignore"} {
				content, ok := out[key]
				if !ok {
					t.Fatalf("missing %s for framework %s", key, fw)
				}
				if len(content) == 0 {
					t.Fatalf("empty %s for framework %s", key, fw)
				}
			}
		})
	}
}

func TestRenderFrameworkPortsAndCompose(t *testing.T) {
	expectedPorts := map[detect.Framework]string{
		detect.Nextjs:     "3000:3000",
		detect.Nestjs:     "3000:3000",
		detect.Express:    "3000:3000",
		detect.React:      "80:80",
		detect.Go:         "8080:8080",
		detect.Rust:       "8080:8080",
		detect.SpringBoot: "8080:8080",
	}

	for fw, wantPort := range expectedPorts {
		t.Run(string(fw), func(t *testing.T) {
			out, err := Render(fw)
			if err != nil {
				t.Fatalf("unexpected error rendering %s: %v", fw, err)
			}

			compose := out["compose.yaml"]
			expectedLine := "- \"" + wantPort + "\""
			if !strings.Contains(compose, expectedLine) {
				t.Errorf("compose.yaml for %s missing port %q, got:\n%s", fw, expectedLine, compose)
			}
			if strings.Contains(compose, "healthcheck:") || strings.Contains(compose, "curl") {
				t.Errorf("compose.yaml for %s should not contain healthcheck block, got:\n%s", fw, compose)
			}
		})
	}
}

func TestSpringBootDockerfileGradleAndMavenSupport(t *testing.T) {
	out, err := Render(detect.SpringBoot)
	if err != nil {
		t.Fatalf("unexpected error rendering SpringBoot: %v", err)
	}
	df := out["Dockerfile"]
	if !strings.Contains(df, "/app/dist/app.jar") {
		t.Errorf("SpringBoot Dockerfile missing /app/dist/app.jar location, got:\n%s", df)
	}
}

func TestRenderDockerfilesAreProductionGrade(t *testing.T) {
	frameworks := []detect.Framework{
		detect.Nextjs,
		detect.React,
		detect.Nestjs,
		detect.Express,
		detect.SpringBoot,
		detect.Go,
		detect.Rust,
	}

	for _, fw := range frameworks {
		t.Run(string(fw), func(t *testing.T) {
			out, err := Render(fw)
			if err != nil {
				t.Fatalf("unexpected error rendering %s: %v", fw, err)
			}

			df := out["Dockerfile"]
			for _, want := range []string{"AS build", "USER ", "HEALTHCHECK"} {
				if !strings.Contains(df, want) {
					t.Fatalf("Dockerfile for %s missing %q:\n%s", fw, want, df)
				}
			}

			if strings.Contains(df, ":latest") {
				t.Fatalf("Dockerfile for %s must not use :latest tag:\n%s", fw, df)
			}
		})
	}
}

func TestRenderUnknownFrameworkErrors(t *testing.T) {
	if _, err := Render(detect.Framework("cobol")); err == nil {
		t.Fatal("expected error for unknown framework")
	}
}
