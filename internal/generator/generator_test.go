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
