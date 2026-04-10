package cli

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	checkStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Bold(true)
	uncheckStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

type multiSelectModel struct {
	prompt   string
	items    []string
	checked  []bool
	cursor   int
	done     bool
	quitting bool
}

func newMultiSelectModel(prompt string, items []string) multiSelectModel {
	return multiSelectModel{
		prompt:  prompt,
		items:   items,
		checked: make([]bool, len(items)),
	}
}

func (m multiSelectModel) Init() tea.Cmd { return nil }

func (m multiSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "ctrl+c", "esc":
			m.quitting = true
			m.done = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case " ", "x":
			m.checked[m.cursor] = !m.checked[m.cursor]
		case "a":
			allChecked := true
			for _, c := range m.checked {
				if !c {
					allChecked = false
					break
				}
			}
			for i := range m.checked {
				m.checked[i] = !allChecked
			}
		case "enter":
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m multiSelectModel) View() string {
	if m.done && !m.quitting {
		var selected []string
		for i, item := range m.items {
			if m.checked[i] {
				selected = append(selected, item)
			}
		}
		if len(selected) == 0 {
			return confirmPromptStyle.Render(m.prompt) + " " + dimStyle.Render("(none)") + "\n"
		}
		return confirmPromptStyle.Render(m.prompt) + " " + pickedStyle.Render(strings.Join(selected, ", ")) + "\n"
	}
	if m.done {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(confirmPromptStyle.Render(m.prompt))
	sb.WriteString(" ")
	sb.WriteString(dimStyle.Render("(space to toggle, a to all, enter to confirm)"))
	sb.WriteString("\n")

	for i, item := range m.items {
		check := uncheckStyle.Render("[ ]")
		if m.checked[i] {
			check = checkStyle.Render("[x]")
		}

		if i == m.cursor {
			sb.WriteString(cursorStyle.Render("  > "))
			sb.WriteString(check)
			sb.WriteString(" ")
			sb.WriteString(selectedStyle.Render(item))
		} else {
			sb.WriteString("    ")
			sb.WriteString(check)
			sb.WriteString(" ")
			sb.WriteString(dimStyle.Render(item))
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (m multiSelectModel) selected() []string {
	var result []string
	for i, item := range m.items {
		if m.checked[i] {
			result = append(result, item)
		}
	}
	return result
}

func (m multiSelectModel) selectedIndices() []int {
	var result []int
	for i, c := range m.checked {
		if c {
			result = append(result, i)
		}
	}
	return result
}

// MultiSelect presents a list of options with checkboxes and returns the selected items.
// Uses space/x to toggle, a to toggle all, arrow keys or j/k for navigation, Enter to confirm.
// Returns an error if the user cancels with Esc or Ctrl+C.
func MultiSelect(prompt string, items []string) ([]string, []int, error) {
	if len(items) == 0 {
		return nil, nil, fmt.Errorf("multiselect: no items provided")
	}
	m := newMultiSelectModel(prompt, items)
	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return nil, nil, fmt.Errorf("multiselect: %w", err)
	}
	final, ok := result.(multiSelectModel)
	if !ok {
		return nil, nil, fmt.Errorf("multiselect: unexpected model type")
	}
	if final.quitting {
		return nil, nil, fmt.Errorf("multiselect: cancelled")
	}
	return final.selected(), final.selectedIndices(), nil
}
