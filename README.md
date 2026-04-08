# tuikit-go

The pragmatic TUI toolkit for shipping CLI tools fast. Wraps [Bubble Tea](https://github.com/charmbracelet/bubbletea) + [Lip Gloss](https://github.com/charmbracelet/lipgloss) with reusable components, a layout engine, a keybinding registry, a theme system, and built-in binary self-update. Build a complete TUI app in under 20 lines.

## Features

- Table with sorting, filtering, custom cell rendering, and mouse support
- ListView, StatusBar, Help screen, ConfigEditor, CommandBar, DetailOverlay, CollapsibleSection
- Dual-pane layout engine with collapsible sidebar
- Keybinding registry with auto-generated help screen
- Dark and light themes with semantic color tokens
- Poller for background data with automatic tick-driven refresh
- Utilities: Sparkline, RelativeTime, OpenURL, OSC8Link
- **Self-update built in** — binary replacement with SHA256 checksum verification, update cache, and Homebrew/Scoop detection. No other Go TUI library ships this.

## Install

```bash
go get github.com/moneycaringcoder/tuikit-go
```

## Quick Start

```go
package main

import (
    "fmt"
    tuikit "github.com/moneycaringcoder/tuikit-go"
)

func main() {
    table := tuikit.NewTable(
        []tuikit.Column{
            {Title: "Name", Width: 20, Sortable: true},
            {Title: "Status", Width: 15},
        },
        []tuikit.Row{
            {"Alice", "Online"},
            {"Bob", "Away"},
        },
        tuikit.TableOpts{Sortable: true, Filterable: true},
    )

    app := tuikit.NewApp(
        tuikit.WithTheme(tuikit.DefaultTheme()),
        tuikit.WithComponent("main", table),
        tuikit.WithStatusBar(
            func() string { return " ? help  q quit" },
            func() string { return fmt.Sprintf(" %d items", 2) },
        ),
        tuikit.WithHelp(),
    )

    app.Run()
}
```

See [`examples/dashboard/`](examples/dashboard/) for a complete app showing all components together.

```bash
go run ./examples/dashboard/
```

## Components

### Table

Adaptive table with responsive columns, sorting, filtering, and cursor navigation.

```go
columns := []tuikit.Column{
    {Title: "Name", Width: 20, Sortable: true},
    {Title: "Score", Width: 10, Align: tuikit.Right, Sortable: true},
    {Title: "Extra", Width: 15, MinWidth: 100}, // hides below 100 cols
}

table := tuikit.NewTable(columns, rows, tuikit.TableOpts{
    Sortable:   true, // 's' to cycle sort
    Filterable: true, // '/' to search
})

table.SetRows(newRows) // update data dynamically
```

**Custom cell rendering** — full control over per-cell styling:

```go
tuikit.TableOpts{
    CellRenderer: func(row tuikit.Row, colIdx int, isCursor bool, theme tuikit.Theme) string {
        val := row[colIdx]
        if colIdx == 2 && val == "Online" {
            return lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Positive)).Render(val)
        }
        return val
    },
}
```

**Custom sort** — numeric, time-based, or any logic:

```go
tuikit.TableOpts{
    SortFunc: func(a, b tuikit.Row, sortCol int, sortAsc bool) bool {
        va, _ := strconv.ParseFloat(a[sortCol], 64)
        vb, _ := strconv.ParseFloat(b[sortCol], 64)
        if sortAsc { return va < vb }
        return va > vb
    },
}
```

**Predicate filter** — programmatic row filtering alongside text search:

```go
table.SetFilter(func(row tuikit.Row) bool {
    return row[1] == "online"
})
table.SetFilter(nil) // clear
```

Mouse scroll and click are handled automatically when `tuikit.WithMouseSupport()` is set.

### StatusBar

Footer with left-aligned hints and right-aligned status.

```go
tuikit.WithStatusBar(
    func() string { return " ? help  q quit" },
    func() string { return " 42 items" },
)
```

### Help

Auto-generated from all registered keybindings. Zero configuration — press `?` to toggle.

```go
tuikit.WithHelp()
```

### ConfigEditor

Declarative settings overlay with grouped fields and validation.

```go
editor := tuikit.NewConfigEditor([]tuikit.ConfigField{
    {
        Label: "Refresh Interval",
        Group: "General",
        Hint:  "seconds, min 5",
        Get:   func() string { return fmt.Sprint(cfg.Interval) },
        Set: func(v string) error {
            n, _ := strconv.Atoi(v)
            if n < 5 { return fmt.Errorf("must be >= 5") }
            cfg.Interval = n
            return nil
        },
    },
})

tuikit.WithOverlay("Settings", "c", editor) // press 'c' to open
```

### CommandBar

Inline command input with completion and dispatch.

### DetailOverlay

Full-screen detail view for a selected row or item.

### CollapsibleSection

Expandable section for grouping content in a panel.

### Poller

Background data polling with tick-driven refresh.

### App-Level Keybindings

```go
tuikit.WithKeyBind(tuikit.KeyBind{
    Key:   "f",
    Label: "Cycle filter",
    Group: "DATA",
    Handler: func() {
        filterIdx = (filterIdx + 1) % len(modes)
        table.SetRows(rows)
    },
})
```

Registered keybindings appear in the help screen automatically.

### Tick / Timer

```go
tuikit.WithTickInterval(100 * time.Millisecond)
```

Components receive `tuikit.TickMsg` in their `Update` method.

### External Data

Push data into the app from WebSocket streams, API polling, or any goroutine:

```go
app := tuikit.NewApp(...)
go func() {
    for data := range stream {
        app.Send(MyDataMsg{data})
    }
}()
app.Run()
```

Unknown message types are forwarded to all components via `Update`.

## Utilities

### Sparkline

Unicode block sparkline from a `[]float64` slice. Bars are colored by direction (up/down/neutral) or rendered mono.

```go
line, width := tuikit.Sparkline(prices, 40, &tuikit.SparklineOpts{Mono: false})
```

### RelativeTime

Short human-readable duration string.

```go
tuikit.RelativeTime(event.Time, time.Now()) // "3m ago", "2h ago", "5d ago"
```

### OpenURL

Opens a URL in the user's default browser. Runs asynchronously, does not block.

```go
tuikit.OpenURL("https://github.com/moneycaringcoder/tuikit-go")
```

### OSC8Link

Wraps text in an OSC8 terminal hyperlink escape sequence.

```go
tuikit.OSC8Link("https://example.com", "click here")
```

## Layout & Theming

### Layout

Single pane or dual pane with collapsible sidebar.

```go
tuikit.WithLayout(&tuikit.DualPane{
    Main:         table,
    Side:         panel,
    SideWidth:    30,
    MinMainWidth: 60, // sidebar auto-hides below this
    SideRight:    true,
    ToggleKey:    "p",
})
```

### Theming

Built-in dark and light themes, or create your own from a color map.

```go
tuikit.DefaultTheme()
tuikit.LightTheme()

// From parsed config (YAML/JSON/TOML)
tuikit.ThemeFromMap(map[string]string{
    "positive": "#00ff00",
    "negative": "#ff0000",
    "accent":   "#0000ff",
})
```

Semantic tokens: `Positive`, `Negative`, `Accent`, `Muted`, `Text`, `TextInverse`, `Cursor`, `Border`, `Flash`.

Components receive theme updates automatically if they implement `Themed`:

```go
type Themed interface {
    SetTheme(tuikit.Theme)
}
```

## Self-Update

tuikit-go ships a complete binary self-update system. No other Go TUI library includes this. It checks GitHub Releases, verifies SHA256 checksums against GoReleaser's `checksums.txt`, replaces the running binary atomically, and detects Homebrew/Scoop installs to skip auto-replace when the package manager owns the binary.

### Built-in app integration

Add one option to your `NewApp` call:

```go
app := tuikit.NewApp(
    tuikit.WithAutoUpdate(tuikit.UpdateConfig{
        Owner:      "myorg",
        Repo:       "mytool",
        BinaryName: "mytool",
        Version:    version, // set via ldflags: -X main.version=v1.2.3
        Mode:       tuikit.UpdateNotify,   // or UpdateBlocking
        CacheTTL:   24 * time.Hour,
    }),
)
```

- `UpdateNotify` — shows a non-blocking banner inside the TUI after startup.
- `UpdateBlocking` — prompts in stdout before the TUI starts.
- Dev builds (`version == ""` or `"dev"`) are skipped automatically.
- Results are cached to avoid hitting the GitHub API on every launch.

### Manual update command

Call `SelfUpdate` directly to implement an explicit `--update` flag or menu action:

```go
if err := tuikit.SelfUpdate(cfg); err != nil {
    fmt.Fprintln(os.Stderr, "update failed:", err)
    os.Exit(1)
}
```

`SelfUpdate` downloads the matching release asset, verifies its checksum, extracts the binary, and replaces the running executable. Add `CleanupOldBinary()` near the top of `main()` to remove the `.old` backup left by a previous update.

### Install method detection

```go
method := tuikit.DetectInstallMethod(os.Args[0])
// InstallManual, InstallHomebrew, or InstallScoop
```

Use this to skip `SelfUpdate` when the binary is managed by a package manager.

## Predictable API

tuikit follows consistent patterns throughout:

- All components implement `tuikit.Component` — `Init`, `Update`, `View`, `KeyBindings`, `SetSize`, `Focused`, `SetFocused`
- App is configured via functional options: `WithTheme`, `WithComponent`, `WithLayout`, `WithAutoUpdate`, etc.
- Key dispatch order: overlay stack → built-in globals (`q` / `tab` / `?`) → pane toggle → overlay triggers → app keybindings → focused component
- Return `tuikit.Consumed()` from a component's `Update` to stop further dispatch
- Modal overlays implement the `Overlay` interface and register with a trigger key via `WithOverlay`
- See `examples/dashboard/main.go` for a complete reference

## Used By

- [cryptstream-tui](https://github.com/moneycaringcoder/cryptstream-tui) — Live cryptocurrency ticker
- [gitstream-tui](https://github.com/moneycaringcoder/gitstream-tui) — GitHub activity dashboard

## Dependencies

Charm ecosystem only:

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Styling
- [Bubbles](https://github.com/charmbracelet/bubbles) — Component primitives

## License

MIT
