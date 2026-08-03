package doctor

import (
	"os"
	"path/filepath"
	"testing"
)

func TestScoreDeductsFailingWeights(t *testing.T) {
	checks := []Check{
		{ID: "a", Severity: Error, Weight: 20, Eval: func(State) bool { return false }},
		{ID: "b", Severity: Warn, Weight: 5, Eval: func(State) bool { return true }},
	}
	rep := Run(State{}, checks)
	if rep.Score != 80 {
		t.Fatalf("want 80, got %d", rep.Score)
	}
}

func TestScoreClampsAtZero(t *testing.T) {
	checks := []Check{
		{ID: "a", Weight: 200, Eval: func(State) bool { return false }},
	}
	if rep := Run(State{}, checks); rep.Score != 0 {
		t.Fatalf("want 0, got %d", rep.Score)
	}
}

func TestDefaultChecksNonEmpty(t *testing.T) {
	if len(DefaultChecks()) == 0 {
		t.Fatal("expected default checks")
	}
}

func TestCollectState(t *testing.T) {
	dir := t.TempDir()

	// Write mock files
	if err := os.WriteFile(filepath.Join(dir, "Dockerfile"), []byte("FROM node:18\nUSER node\nHEALTHCHECK CMD exit 0"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "compose.yaml"), []byte("services:\n  app:\n    restart: always"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".dockerignore"), []byte("node_modules"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("PORT=3000"), 0644); err != nil {
		t.Fatal(err)
	}

	state := CollectState(dir)
	if !state.HasDockerfile {
		t.Error("expected HasDockerfile to be true")
	}
	if !state.HasCompose {
		t.Error("expected HasCompose to be true")
	}
	if !state.HasDockerignore {
		t.Error("expected HasDockerignore to be true")
	}
	if !state.HasEnvFile {
		t.Error("expected HasEnvFile to be true")
	}
	if !state.HasHealthcheck {
		t.Error("expected HasHealthcheck to be true")
	}
	if !state.HasRestart {
		t.Error("expected HasRestart to be true")
	}
	if state.RunsAsRoot {
		t.Error("expected RunsAsRoot to be false because USER node is set")
	}
	if state.UsesLatest {
		t.Error("expected UsesLatest to be false")
	}

	report := Run(state, DefaultChecks())
	if report.Score == 0 {
		t.Error("expected non-zero score for valid project")
	}
}

func TestCollectStateDetectsPortConflictAndMissingVolume(t *testing.T) {
	dir := t.TempDir()
	compose := `services:
  api:
    ports:
      - "3000:3000"
    volumes:
      - db_data:/var/lib/data
  worker:
    ports:
      - "3000:9000"
volumes: {}
`
	if err := os.WriteFile(filepath.Join(dir, "compose.yaml"), []byte(compose), 0644); err != nil {
		t.Fatal(err)
	}
	state := CollectState(dir)
	if len(state.PortConflicts) != 1 || state.PortConflicts[0] != "3000" {
		t.Fatalf("want port conflict on 3000, got %v", state.PortConflicts)
	}
	if len(state.MissingVolumes) != 1 || state.MissingVolumes[0] != "db_data" {
		t.Fatalf("want missing volume db_data, got %v", state.MissingVolumes)
	}
}

func TestCollectStateCleanComposeHasNoConflicts(t *testing.T) {
	dir := t.TempDir()
	compose := `services:
  api:
    ports:
      - "3000:3000"
    volumes:
      - db_data:/var/lib/data
      - ./src:/app
volumes:
  db_data:
`
	if err := os.WriteFile(filepath.Join(dir, "compose.yaml"), []byte(compose), 0644); err != nil {
		t.Fatal(err)
	}
	state := CollectState(dir)
	if len(state.PortConflicts) != 0 {
		t.Fatalf("expected no port conflicts, got %v", state.PortConflicts)
	}
	if len(state.MissingVolumes) != 0 {
		t.Fatalf("expected no missing volumes (bind mount ./src ignored), got %v", state.MissingVolumes)
	}
}
