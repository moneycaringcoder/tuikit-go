package cli

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	passwordDoneStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Bold(true)
)

type passwordModel struct {
	prompt   string
	validate func(string) error
	ti       textinput.Model
	errMsg   string
	value    string
	done     bool
	quitting bool
}

func newPasswordModel(prompt string, validate func(string) error) passwordModel {
	ti := textinput.New()
	ti.Focus()
	ti.Prompt = ""
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'
	return passwordModel{
		prompt:   prompt,
		validate: validate,
		ti:       ti,
	}
}

func (m passwordModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m passwordModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "ctrl+c":
			m.quitting = true
			m.done = true
			return m, tea.Quit
		case "enter":
			val := m.ti.Value()
			if m.validate != nil {
				if err := m.validate(val); err != nil {
					m.errMsg = err.Error()
					m.ti.SetValue("")
					return m, nil
				}
			}
			m.value = val
			m.done = true
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.ti, cmd = m.ti.Update(msg)
	return m, cmd
}

func (m passwordModel) View() string {
	if m.done && !m.quitting {
		masked := strings.Repeat("•", len(m.value))
		return inputPromptStyle.Render(m.prompt) + " " + passwordDoneStyle.Render(masked) + "\n"
	}
	if m.done {
		return ""
	}
	var sb strings.Builder
	sb.WriteString(inputPromptStyle.Render(m.prompt))
	sb.WriteString(" ")
	sb.WriteString(m.ti.View())
	if m.errMsg != "" {
		sb.WriteString("\n  ")
		sb.WriteString(inputErrStyle.Render(m.errMsg))
	}
	return sb.String()
}

// Password prompts for masked text input with optional validation.
// Characters are displayed as bullet points (•).
// Returns an error only if the user cancels with Ctrl+C.
func Password(prompt string, validate func(string) error) (string, error) {
	m := newPasswordModel(prompt, validate)
	p := tea.NewProgram(m)
	result, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("password: %w", err)
	}
	final, ok := result.(passwordModel)
	if !ok {
		return "", fmt.Errorf("password: unexpected model type")
	}
	if final.quitting {
		return "", fmt.Errorf("password: cancelled")
	}
	return final.value, nil
}
