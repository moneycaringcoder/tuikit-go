package cli

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
	successIcon  = lipgloss.NewStyle().Foreground(lipgloss.Color("82")).Bold(true).Render("✓")
	warningIcon  = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Bold(true).Render("!")
	errorIcon    = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true).Render("✗")
	infoIcon     = lipgloss.NewStyle().Foreground(lipgloss.Color("75")).Bold(true).Render("ℹ")
	successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
	warningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	infoStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("75"))
	stepStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	stepNumStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	titleStyle   = lipgloss.NewStyle().Bold(true).Underline(true)
	sectionStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
)

// Success prints a green success message with a checkmark.
func Success(msg string) {
	fmt.Printf("  %s %s\n", successIcon, successStyle.Render(msg))
}

// Warning prints a yellow warning message with an exclamation mark.
func Warning(msg string) {
	fmt.Printf("  %s %s\n", warningIcon, warningStyle.Render(msg))
}

// Error prints a red error message with an X mark.
func Error(msg string) {
	fmt.Printf("  %s %s\n", errorIcon, errorStyle.Render(msg))
}

// Info prints a blue informational message with an info icon.
func Info(msg string) {
	fmt.Printf("  %s %s\n", infoIcon, infoStyle.Render(msg))
}

// Successf prints a formatted green success message.
func Successf(format string, a ...any) {
	Success(fmt.Sprintf(format, a...))
}

// Warningf prints a formatted yellow warning message.
func Warningf(format string, a ...any) {
	Warning(fmt.Sprintf(format, a...))
}

// Errorf prints a formatted red error message.
func Errorf(format string, a ...any) {
	Error(fmt.Sprintf(format, a...))
}

// Infof prints a formatted blue informational message.
func Infof(format string, a ...any) {
	Info(fmt.Sprintf(format, a...))
}

// Step prints a numbered step indicator. Useful for multi-step CLI flows.
// Example: Step(1, 3, "Installing dependencies")  →  "  [1/3] Installing dependencies"
func Step(current, total int, msg string) {
	num := stepNumStyle.Render(fmt.Sprintf("[%d/%d]", current, total))
	fmt.Printf("  %s %s\n", num, stepStyle.Render(msg))
}

// Title prints a bold underlined title.
func Title(msg string) {
	fmt.Printf("\n  %s\n\n", titleStyle.Render(msg))
}

// Section prints a bold colored section header.
func Section(msg string) {
	fmt.Printf("\n  %s\n", sectionStyle.Render(msg))
}

// Separator prints a horizontal divider line.
func Separator() {
	fmt.Println(dimStyle.Render("  ────────────────────────────────────────"))
}

// Dim prints dimmed/muted text.
func Dim(msg string) {
	fmt.Printf("  %s\n", dimStyle.Render(msg))
}

// KeyValue prints a key-value pair with the key dimmed and value normal.
func KeyValue(key, value string) {
	fmt.Printf("  %s %s\n", dimStyle.Render(key+":"), value)
}
