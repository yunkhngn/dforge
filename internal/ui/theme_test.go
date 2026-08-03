package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestSymbols(t *testing.T) {
	if OK != "✓" || WARN != "⚠" || ERR != "✗" {
		t.Fatalf("unexpected symbols: %q %q %q", OK, WARN, ERR)
	}
}

func TestSuccessContainsText(t *testing.T) {
	if got := Success("done"); got == "" {
		t.Fatal("Success returned empty")
	}
}

func TestWarningContainsText(t *testing.T) {
	if got := Warning("warn"); got == "" {
		t.Fatal("Warning returned empty")
	}
}

func TestErrorContainsText(t *testing.T) {
	if got := Error("err"); got == "" {
		t.Fatal("Error returned empty")
	}
}

func TestTitleContainsText(t *testing.T) {
	if got := Title("title"); got == "" {
		t.Fatal("Title returned empty")
	}
}

func TestSelectModel(t *testing.T) {
	options := []string{"Option 1", "Option 2", "Option 3"}
	m := selectModel{
		prompt:  "Select an option",
		options: options,
		chosen:  -1,
	}

	if m.Init() != nil {
		t.Errorf("Init() should return nil")
	}

	view := m.View()
	if !strings.Contains(view, "Select an option") || !strings.Contains(view, "> Option 1") {
		t.Errorf("View output unexpected: %s", view)
	}

	// Move down
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(selectModel)
	if cmd != nil {
		t.Errorf("KeyDown should not return cmd")
	}
	if m.cursor != 1 {
		t.Errorf("expected cursor 1, got %d", m.cursor)
	}

	// Move down again
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(selectModel)
	if m.cursor != 2 {
		t.Errorf("expected cursor 2, got %d", m.cursor)
	}

	// Try move down past end
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(selectModel)
	if m.cursor != 2 {
		t.Errorf("expected cursor to stay 2, got %d", m.cursor)
	}

	// Move up
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(selectModel)
	if m.cursor != 1 {
		t.Errorf("expected cursor 1, got %d", m.cursor)
	}

	// Press enter
	updated, cmd = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(selectModel)
	if m.chosen != 1 || !m.done {
		t.Errorf("expected chosen 1 and done true, got chosen %d, done %v", m.chosen, m.done)
	}
	if cmd == nil {
		t.Errorf("Enter should return quit cmd")
	}

	// Test cancellation with q
	m2 := selectModel{prompt: "Test", options: options, chosen: -1}
	updated, _ = m2.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
	m2 = updated.(selectModel)
	if m2.chosen != -1 || !m2.done {
		t.Errorf("expected chosen -1 and done true on q, got chosen %d, done %v", m2.chosen, m2.done)
	}
}

func TestConfirmWithReader(t *testing.T) {
	tests := []struct {
		input      string
		defaultYes bool
		expected   bool
	}{
		{"\n", true, true},
		{"\n", false, false},
		{"y\n", false, true},
		{"Y\n", false, true},
		{"yes\n", false, true},
		{"n\n", true, false},
		{"N\n", true, false},
		{"no\n", true, false},
	}

	for _, tt := range tests {
		reader := strings.NewReader(tt.input)
		got, err := confirmWithReader(reader, "Continue?", tt.defaultYes)
		if err != nil {
			t.Errorf("confirmWithReader(%q, %v) unexpected error: %v", tt.input, tt.defaultYes, err)
		}
		if got != tt.expected {
			t.Errorf("confirmWithReader(%q, %v) = %v; want %v", tt.input, tt.defaultYes, got, tt.expected)
		}
	}
}
