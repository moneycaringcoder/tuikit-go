# Quick Start

## Install

```bash
go get github.com/moneycaringcoder/tuikit-go
```

Requires Go 1.24+.

## Minimal App

The smallest possible tuikit-go app registers one component and runs:

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

Run it:

```bash
go run .
```

Keys: `j`/`k` to move, `s` to cycle sort, `/` to search, `?` for help, `q` to quit.

## What `NewApp` wires up

| Option | Effect |
|--------|--------|
| `WithTheme` | Applies semantic color tokens to all components |
| `WithComponent` | Registers a component as the main pane |
| `WithLayout` | Dual-pane layout with sidebar |
| `WithStatusBar` | Footer with left/right text |
| `WithHelp` | `?` toggle overlay — auto-populated from all `KeyBindings()` |
| `WithMouseSupport` | Enables mouse scroll and click |
| `WithAutoUpdate` | Binary self-update on launch |

## Next Steps

- [App Structure](app-structure.md) — component interface, slots, key dispatch
- [Theming](theming.md) — dark/light themes, custom tokens
- [Testing](testing.md) — tuitest virtual terminal assertions
