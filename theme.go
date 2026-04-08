package tuikit

import "github.com/charmbracelet/lipgloss"

// Theme defines semantic color tokens for consistent styling across components.
// Components reference these tokens instead of raw colors.
type Theme struct {
	Positive    lipgloss.Color // Green: gains, success, online
	Negative    lipgloss.Color // Red: losses, errors, offline
	Accent      lipgloss.Color // Highlights, active elements
	Muted       lipgloss.Color // Dimmed text, secondary info
	Text        lipgloss.Color // Primary text
	TextInverse lipgloss.Color // Text on colored backgrounds
	Cursor      lipgloss.Color // Cursor/selection highlight
	Border      lipgloss.Color // Borders, separators
	Flash       lipgloss.Color // Temporary notification background
	Extra       map[string]lipgloss.Color // App-specific color tokens
}

// Color returns an app-specific color from Extra by key, falling back to
// the provided default if the key doesn't exist.
func (t Theme) Color(key string, fallback lipgloss.Color) lipgloss.Color {
	if t.Extra != nil {
		if c, ok := t.Extra[key]; ok {
			return c
		}
	}
	return fallback
}

// DefaultTheme returns a dark theme suitable for most terminal backgrounds.
func DefaultTheme() Theme {
	return Theme{
		Positive:    lipgloss.Color("#22c55e"),
		Negative:    lipgloss.Color("#ef4444"),
		Accent:      lipgloss.Color("#3b82f6"),
		Muted:       lipgloss.Color("#6b7280"),
		Text:        lipgloss.Color("#e5e7eb"),
		TextInverse: lipgloss.Color("#111827"),
		Cursor:      lipgloss.Color("#38bdf8"),
		Border:      lipgloss.Color("#374151"),
		Flash:       lipgloss.Color("#facc15"),
	}
}

// LightTheme returns a light theme for light terminal backgrounds.
func LightTheme() Theme {
	return Theme{
		Positive:    lipgloss.Color("#16a34a"),
		Negative:    lipgloss.Color("#dc2626"),
		Accent:      lipgloss.Color("#2563eb"),
		Muted:       lipgloss.Color("#9ca3af"),
		Text:        lipgloss.Color("#111827"),
		TextInverse: lipgloss.Color("#f9fafb"),
		Cursor:      lipgloss.Color("#0284c7"),
		Border:      lipgloss.Color("#d1d5db"),
		Flash:       lipgloss.Color("#eab308"),
	}
}

// ThemeFromMap creates a Theme from a map of color names to hex values.
// Missing keys fall back to DefaultTheme values. This is config-format-agnostic:
// your app reads YAML/JSON/TOML and passes the color map here.
// Keys not matching a built-in token are placed in the Extra map.
func ThemeFromMap(m map[string]string) Theme {
	t := DefaultTheme()
	builtins := map[string]*lipgloss.Color{
		"positive":     &t.Positive,
		"negative":     &t.Negative,
		"accent":       &t.Accent,
		"muted":        &t.Muted,
		"text":         &t.Text,
		"text_inverse": &t.TextInverse,
		"cursor":       &t.Cursor,
		"border":       &t.Border,
		"flash":        &t.Flash,
	}
	for k, v := range m {
		if ptr, ok := builtins[k]; ok {
			*ptr = lipgloss.Color(v)
		} else {
			if t.Extra == nil {
				t.Extra = make(map[string]lipgloss.Color)
			}
			t.Extra[k] = lipgloss.Color(v)
		}
	}
	return t
}
