package ui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type selectModel struct {
	prompt  string
	options []string
	cursor  int
	chosen  int
	done    bool
}

func (m selectModel) Init() tea.Cmd { return nil }

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.options)-1 {
				m.cursor++
			}
		case "enter":
			m.chosen = m.cursor
			m.done = true
			return m, tea.Quit
		case "ctrl+c", "q":
			m.chosen = -1
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m selectModel) View() string {
	s := Title(m.prompt) + "\n\n"
	for i, opt := range m.options {
		cursor := "  "
		if i == m.cursor {
			cursor = "> "
		}
		s += cursor + opt + "\n"
	}
	return s
}

// Select runs an interactive single-select picker and returns the chosen index.
func Select(prompt string, options []string) (int, error) {
	m := selectModel{prompt: prompt, options: options, chosen: -1}
	res, err := tea.NewProgram(m).Run()
	if err != nil {
		return -1, err
	}
	final := res.(selectModel)
	if final.chosen < 0 {
		return -1, fmt.Errorf("selection cancelled")
	}
	return final.chosen, nil
}
