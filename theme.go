package tuikit

import (
	"math"
	"strconv"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Theme defines semantic color tokens for consistent styling across components.
// Components reference these tokens instead of raw colors.
type Theme struct {
	Positive    lipgloss.Color
	Negative    lipgloss.Color
	Accent      lipgloss.Color
	Muted       lipgloss.Color
	Text        lipgloss.Color
	TextInverse lipgloss.Color
	Cursor      lipgloss.Color
	Border      lipgloss.Color
	Flash       lipgloss.Color
	Extra       map[string]lipgloss.Color

	// Glyphs holds the symbol set for this theme.
	// If nil, components fall back to DefaultGlyphs().
	Glyphs *Glyphs

	// Borders holds named border styles for this theme.
	// If nil, components fall back to DefaultBorders().
	Borders *BorderSet
}

// glyphsOrDefault returns the theme's Glyphs, falling back to DefaultGlyphs.
func (t Theme) glyphsOrDefault() Glyphs {
	if t.Glyphs != nil {
		return *t.Glyphs
	}
	return DefaultGlyphs()
}

// bordersOrDefault returns the theme's BorderSet, falling back to DefaultBorders.
func (t Theme) bordersOrDefault() BorderSet {
	if t.Borders != nil {
		return *t.Borders
	}
	return DefaultBorders()
}

// BorderSet holds named lipgloss border styles for a theme.
type BorderSet struct {
	Rounded lipgloss.Border
	Double  lipgloss.Border
	Thick   lipgloss.Border
	Ascii   lipgloss.Border
	Minimal lipgloss.Border
}

// DefaultBorders returns standard lipgloss border presets.
func DefaultBorders() BorderSet {
	return BorderSet{
		Rounded: lipgloss.RoundedBorder(),
		Double:  lipgloss.DoubleBorder(),
		Thick:   lipgloss.ThickBorder(),
		Ascii:   lipgloss.ASCIIBorder(),
		Minimal: lipgloss.NormalBorder(),
	}
}

// --- theme registry ---

var (
	themeMu       sync.RWMutex
	themeRegistry = map[string]Theme{}
)

// Register adds a named theme to the global registry.
// Calling Register with the same name twice overwrites the previous entry.
func Register(name string, t Theme) {
	themeMu.Lock()
	defer themeMu.Unlock()
	themeRegistry[name] = t
}

// Presets returns a copy of all registered themes keyed by name.
func Presets() map[string]Theme {
	themeMu.RLock()
	defer themeMu.RUnlock()
	out := make(map[string]Theme, len(themeRegistry))
	for k, v := range themeRegistry {
		out[k] = v
	}
	return out
}

// Contrast returns the WCAG relative contrast ratio between two hex colours.
// A value >= 4.5 satisfies WCAG AA for normal text.
func Contrast(bg, fg lipgloss.Color) float64 {
	l1 := relativeLuminance(string(bg))
	l2 := relativeLuminance(string(fg))
	lighter := math.Max(l1, l2)
	darker := math.Min(l1, l2)
	return (lighter + 0.05) / (darker + 0.05)
}

// SetThemeMsg is a tea.Msg that switches the active App theme at runtime.
// Return SetThemeCmd from a WithKeyBind HandlerCmd to change the theme
// for all components immediately.
type SetThemeMsg struct {
	Theme Theme
}

// SetThemeCmd returns a tea.Cmd that sends a SetThemeMsg.
func SetThemeCmd(t Theme) func() tea.Msg {
	return func() tea.Msg { return SetThemeMsg{Theme: t} }
}

// relativeLuminance computes WCAG 2.1 relative luminance for a "#rrggbb" colour.
func relativeLuminance(hex string) float64 {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) != 6 {
		return 0
	}
	rv, _ := strconv.ParseUint(hex[0:2], 16, 8)
	gv, _ := strconv.ParseUint(hex[2:4], 16, 8)
	bv, _ := strconv.ParseUint(hex[4:6], 16, 8)
	lin := func(c uint64) float64 {
		s := float64(c) / 255.0
		if s <= 0.04045 {
			return s / 12.92
		}
		return math.Pow((s+0.055)/1.055, 2.4)
	}
	return 0.2126*lin(rv) + 0.7152*lin(gv) + 0.0722*lin(bv)
}

// Color returns an app-specific color from Extra by key, falling back to
// the provided default if the key does not exist.
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
// Missing keys fall back to DefaultTheme values.
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
