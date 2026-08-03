package cmd

import (
	"strings"
	"testing"

	"github.com/yunkhngn/dforge/internal/env"
)

func TestFormatEnvReportShowsNamesOnly(t *testing.T) {
	r := env.Report{Missing: []string{"SECRET_KEY"}, Empty: []string{"API_TOKEN"}}
	out := formatEnvReport(r)
	if !strings.Contains(out, "SECRET_KEY") || !strings.Contains(out, "API_TOKEN") {
		t.Fatalf("names missing:\n%s", out)
	}
}

func TestFormatEnvReportValid(t *testing.T) {
	r := env.Report{}
	out := formatEnvReport(r)
	if !strings.Contains(out, ".env is valid") {
		t.Fatalf("expected valid message, got:\n%s", out)
	}
}

func TestFormatEnvReportMissingFile(t *testing.T) {
	r := env.Report{MissingFile: true}
	out := formatEnvReport(r)
	if !strings.Contains(out, ".env file not found") {
		t.Fatalf("expected file not found message, got:\n%s", out)
	}
}

func TestFormatEnvReportExtraAndDuplicates(t *testing.T) {
	r := env.Report{
		Extra:      []string{"EXTRA_VAR"},
		Duplicates: []string{"DUP_VAR"},
	}
	out := formatEnvReport(r)
	if !strings.Contains(out, "EXTRA_VAR") || !strings.Contains(out, "DUP_VAR") {
		t.Fatalf("extra/duplicate names missing:\n%s", out)
	}
}
