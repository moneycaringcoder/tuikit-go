# tuikit-go

The pragmatic TUI toolkit for shipping Go CLI tools fast. Wraps [Bubble Tea](https://github.com/charmbracelet/bubbletea) + [Lip Gloss](https://github.com/charmbracelet/lipgloss) with reusable components, a layout engine, a keybinding registry, a theme system, and built-in binary self-update.

```bash
go get github.com/moneycaringcoder/tuikit-go
```

## Features

- **Table** — adaptive columns, sorting, filtering, custom cell rendering, detail bars, virtual mode for millions of rows
- **ListView** — generic scrollable list with cursor, header, detail bar, and flash highlighting
- **LogViewer** — append-only auto-scrolling log with level filtering and substring search
- **Tabs** — horizontal or vertical tab switcher with click-to-focus content panes
- **Picker** — fuzzy-searchable command palette (Cmd-K style) with optional preview pane
- **StatusBar** — left/right footer driven by closures or reactive signals
- **Help** — auto-generated keybinding reference from the registry; press `?` to toggle
- **Dual-pane layout** with collapsible sidebar
- **Dark + light themes** with semantic color tokens and extensible `Extra` color map
- **CLI primitives** — Confirm, SelectOne, MultiSelect, Input, Password, Spinner, Progress
- **tuitest** — virtual-terminal testing framework with 30+ assertions, golden files, and a vitest-style reporter
- **Self-update** — binary replacement with SHA256 verification, rollback, and Homebrew/Scoop detection

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
            {Title: "Name",   Width: 20, Sortable: true},
            {Title: "Status", Width: 15},
        },
        []tuikit.Row{
            {"Alice", "Online"},
            {"Bob",   "Away"},
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

See [Getting Started](guides/quickstart.md) to build your first app in under 5 minutes.

## Examples

```bash
go run ./examples/minimal/     # Simple ListView in ~30 lines
go run ./examples/dashboard/   # Table + DualPane + ConfigEditor
go run ./examples/monitor/     # Service fleet dashboard
go run ./examples/cli-demo/    # Interactive CLI primitives showcase
```

## License

MIT — [github.com/moneycaringcoder/tuikit-go](https://github.com/moneycaringcoder/tuikit-go)
