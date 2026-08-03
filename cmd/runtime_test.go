package cmd

import (
	"testing"

	"github.com/yunkhngn/dforge/internal/docker"
)

func TestRunLogsPassesService(t *testing.T) {
	m := &docker.MockClient{}
	if err := runLogs(m, "api", false); err != nil {
		t.Fatal(err)
	}
	if m.LogsService != "api" {
		t.Fatalf("want api, got %q", m.LogsService)
	}
}

func TestRunCleanRefusesWithoutConfirm(t *testing.T) {
	m := &docker.MockClient{}
	if err := runClean(m, false, false); err == nil {
		t.Fatal("expected refusal without --force/--yes")
	}
	if m.CleanForce {
		t.Fatal("clean should not have run")
	}
}

func TestRunCleanRunsWithForce(t *testing.T) {
	m := &docker.MockClient{}
	if err := runClean(m, true, false); err != nil {
		t.Fatal(err)
	}
	if !m.CleanForce {
		t.Fatal("clean should have run with force")
	}
}
