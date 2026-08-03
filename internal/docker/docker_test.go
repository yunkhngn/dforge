package docker

import (
	"reflect"
	"testing"
)

func TestMockClientSatisfiesInterface(t *testing.T) {
	var c Client = &MockClient{Statuses: []ServiceStatus{{Name: "api", State: "running"}}}
	st, err := c.Status()
	if err != nil {
		t.Fatal(err)
	}
	if len(st) != 1 || st[0].Name != "api" {
		t.Fatalf("unexpected status: %+v", st)
	}
}

func TestMockCleanRecordsForce(t *testing.T) {
	m := &MockClient{}
	if err := m.Clean(true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !m.CleanForce {
		t.Fatal("force not recorded")
	}
}

func TestMockClientMethods(t *testing.T) {
	m := &MockClient{}

	if err := m.Logs("postgres", true); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.LogsService != "postgres" {
		t.Errorf("expected LogsService to be 'postgres', got %q", m.LogsService)
	}

	if err := m.Shell("redis"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.ShellService != "redis" {
		t.Errorf("expected ShellService to be 'redis', got %q", m.ShellService)
	}

	if err := m.Available(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewReturnsClient(t *testing.T) {
	c := New()
	if c == nil {
		t.Fatal("expected New() to return non-nil Client")
	}
}

func TestSplitLines(t *testing.T) {
	input := []byte("line1\nline2\n\nline3")
	expected := [][]byte{
		[]byte("line1"),
		[]byte("line2"),
		[]byte(""),
		[]byte("line3"),
	}

	got := splitLines(input)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("splitLines(%q) = %q, expected %q", string(input), got, expected)
	}
}
