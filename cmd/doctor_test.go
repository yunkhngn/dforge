package cmd

import (
	"strings"
	"testing"

	"github.com/yunkhngn/dforge/internal/doctor"
)

func TestFormatDoctorReportShowsScore(t *testing.T) {
	r := doctor.Report{
		Score: 82,
		Results: []doctor.Result{
			{ID: "docker", Title: "Docker installed", Passed: true, Severity: doctor.Error},
			{ID: "healthcheck", Title: "HEALTHCHECK present", Passed: false, Severity: doctor.Warn},
		},
	}
	out := formatDoctorReport(r)
	if !strings.Contains(out, "Score: 82/100") {
		t.Fatalf("score line missing:\n%s", out)
	}
	if !strings.Contains(out, "Docker installed") || !strings.Contains(out, "HEALTHCHECK present") {
		t.Fatalf("results missing:\n%s", out)
	}
}
