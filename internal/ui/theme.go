package ui

import "github.com/charmbracelet/lipgloss"

const (
	OK   = "✓"
	WARN = "⚠"
	ERR  = "✗"
)

var (
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
	warnStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	errStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	titleStyle   = lipgloss.NewStyle().Bold(true)
)

func Success(s string) string { return successStyle.Render(s) }
func Warning(s string) string { return warnStyle.Render(s) }
func Error(s string) string   { return errStyle.Render(s) }
func Title(s string) string   { return titleStyle.Render(s) }
