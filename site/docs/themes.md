# Themes

tuikit-go ships 8 built-in theme presets plus a `DefaultTheme` (dark teal) and `LightTheme`. All components update automatically when the theme changes.

## Built-In Presets

Use `tuikit.ThemePreset(name)` or call the constructor directly:

```go
theme := tuikit.DraculaTheme()
theme  = tuikit.CatppuccinMochaTheme()
theme  = tuikit.TokyoNightTheme()
theme  = tuikit.NordTheme()
theme  = tuikit.GruvboxDarkTheme()
theme  = tuikit.RosePineTheme()
theme  = tuikit.KanagawaTheme()
theme  = tuikit.OneDarkTheme()
```

Or by name (registered in `init()`):

```go
theme := tuikit.ThemePreset("dracula")
theme  = tuikit.ThemePreset("catppuccin-mocha")
theme  = tuikit.ThemePreset("tokyo-night")
theme  = tuikit.ThemePreset("nord")
theme  = tuikit.ThemePreset("gruvbox-dark")
theme  = tuikit.ThemePreset("rose-pine")
theme  = tuikit.ThemePreset("kanagawa")
theme  = tuikit.ThemePreset("one-dark")
```

## Preset Color Tokens

| Preset | Accent | Cursor | Positive | Negative | Flash |
|--------|--------|--------|----------|----------|-------|
| **Dracula** | `#bd93f9` | `#ff79c6` | `#50fa7b` | `#ff5555` | `#f1fa8c` |
| **Catppuccin Mocha** | `#cba6f7` | `#89b4fa` | `#a6e3a1` | `#f38ba8` | `#f9e2af` |
| **Tokyo Night** | `#7aa2f7` | `#bb9af7` | `#9ece6a` | `#f7768e` | `#e0af68` |
| **Nord** | `#81a1c1` | `#88c0d0` | `#a3be8c` | `#bf616a` | `#ebcb8b` |
| **Gruvbox Dark** | `#fabd2f` | `#83a598` | `#b8bb26` | `#fb4934` | `#d3869b` |
| **Rose Pine** | `#c4a7e7` | `#9ccfd8` | `#31748f` | `#eb6f92` | `#f6c177` |
| **Kanagawa** | `#7e9cd8` | `#957fb8` | `#98bb6c` | `#e46876` | `#e98a00` |
| **One Dark** | `#61afef` | `#c678dd` | `#98c379` | `#e06c75` | `#e5c07b` |

## Default and Light Themes

```go
tuikit.DefaultTheme() // dark teal, accent #14b8a6
tuikit.LightTheme()   // light background with dark text
```

## Building a Custom Theme

Construct a `Theme` directly or parse from a color map:

```go
// Direct construction
theme := tuikit.Theme{
    Positive:    lipgloss.Color("#22c55e"),
    Negative:    lipgloss.Color("#ef4444"),
    Accent:      lipgloss.Color("#3b82f6"),
    Muted:       lipgloss.Color("#6b7280"),
    Text:        lipgloss.Color("#f1f5f9"),
    TextInverse: lipgloss.Color("#0f172a"),
    Cursor:      lipgloss.Color("#2563eb"),
    Border:      lipgloss.Color("#334155"),
    Flash:       lipgloss.Color("#f59e0b"),
}

// From a map (e.g. parsed YAML/TOML config)
theme = tuikit.ThemeFromMap(map[string]string{
    "positive": "#22c55e",
    "negative": "#ef4444",
    "accent":   "#3b82f6",
    // Extra app-specific tokens go into Theme.Extra:
    "pr":     "#3b82f6",
    "review": "#a855f7",
})
```

## Applying a Theme

Pass the theme when creating the app:

```go
app := tuikit.NewApp(
    tuikit.WithTheme(tuikit.DraculaTheme()),
    ...
)
```

Switch at runtime (e.g. from a settings menu):

```go
app.SetTheme(tuikit.TokyoNightTheme())
```

## Hot Reload

Set `TUIKIT_THEME` to a JSON file path. tuikit-go watches the file and reloads the theme without restarting:

```bash
TUIKIT_THEME=~/.config/mytool/theme.json mytool
```

## Custom Glyphs

Override the cursor marker and flash marker glyphs:

```go
theme := tuikit.DefaultTheme()
theme.Glyphs = &tuikit.Glyphs{
    CursorMarker: "▶",
    FlashMarker:  "★",
}
```
