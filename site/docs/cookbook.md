# Cookbook

Five complete recipes for common tuikit-go patterns.

---

## 1. Dashboard in 5 Minutes

A full dashboard with a table, dual-pane sidebar, and status bar:

```go
package main

import (
    "fmt"
    "time"
    tuikit "github.com/moneycaringcoder/tuikit-go"
)

func main() {
    table := tuikit.NewTable(
        []tuikit.Column{
            {Title: "Service", Width: 20, Sortable: true},
            {Title: "Status",  Width: 12},
            {Title: "Latency", Width: 10, Align: tuikit.Right, Sortable: true},
        },
        []tuikit.Row{
            {"api-gateway",  "online",  "12ms"},
            {"auth-service", "online",  "8ms"},
            {"db-primary",   "online",  "2ms"},
            {"cache",        "degraded","45ms"},
        },
        tuikit.TableOpts{
            Sortable:   true,
            Filterable: true,
            CellRenderer: func(row tuikit.Row, col int, isCursor bool, th tuikit.Theme) string {
                if col == 1 {
                    color := th.Positive
                    if row[1] != "online" { color = th.Flash }
                    return lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(row[1])
                }
                return row[col]
            },
        },
    )

    detail := tuikit.NewListView[string](tuikit.ListViewOpts[string]{
        RenderItem: func(s string, _ int, _ bool, th tuikit.Theme) string { return s },
    })

    table.Opts().OnCursorChange = func(row tuikit.Row, _ int) {
        detail.SetItems([]string{
            "Service: " + row[0],
            "Status:  " + row[1],
            "Latency: " + row[2],
        })
    }

    app := tuikit.NewApp(
        tuikit.WithTheme(tuikit.DefaultTheme()),
        tuikit.WithLayout(&tuikit.DualPane{
            Main:         table,
            Side:         detail,
            SideWidth:    28,
            MinMainWidth: 60,
            SideRight:    true,
            ToggleKey:    "p",
        }),
        tuikit.WithStatusBar(
            func() string { return " ? help  s sort  / search  p panel  q quit" },
            func() string { return fmt.Sprintf(" %d services", table.RowCount()) },
        ),
        tuikit.WithHelp(),
        tuikit.WithTickInterval(5*time.Second),
    )
    app.Run()
}
```

---

## 2. Self-Updating CLI

Wire binary self-update into an existing app with three lines:

```go
app := tuikit.NewApp(
    tuikit.WithTheme(tuikit.DefaultTheme()),
    tuikit.WithComponent("main", myMainComponent),
    tuikit.WithAutoUpdate(tuikit.UpdateConfig{
        Owner:      "myorg",
        Repo:       "mytool",
        BinaryName: "mytool",
        Version:    version, // injected via: -ldflags "-X main.version=v1.2.3"
        Mode:       tuikit.UpdateNotify,
        CacheTTL:   24 * time.Hour,
    }),
)
```

Add cleanup at the top of `main()` for the `.old` backup left by a previous update:

```go
func main() {
    tuikit.CleanupOldBinary()
    // ...
}
```

---

## 3. Testing a TUI

Use `tuitest` to assert on rendered screen content without a real terminal:

```go
func TestTable(t *testing.T) {
    table := tuikit.NewTable(
        []tuikit.Column{{Title: "Name", Width: 20}},
        []tuikit.Row{{"Alice"}, {"Bob"}},
        tuikit.TableOpts{Filterable: true},
    )

    tm := tuitest.NewTestModel(t, wrapInApp(table), 80, 24)

    // Header visible
    tuitest.AssertRowContains(t, tm.Screen(), 0, "Name")

    // Navigate down
    tm.SendKey("j")
    tuitest.AssertContains(t, tm.Screen(), "Bob")

    // Filter
    tm.SendKey("/")
    tm.Type("ali")
    tuitest.AssertContains(t, tm.Screen(), "Alice")
    tuitest.AssertNotContains(t, tm.Screen(), "Bob")

    // Golden snapshot
    tuitest.AssertGolden(t, tm.Screen(), "table-filtered")
}
```

Regenerate golden files:

```bash
tuitest -update ./...
```

---

## 4. Importing a Theme

Load a theme from a TOML/YAML config file at startup:

```go
// config.toml:
// [theme]
// accent  = "#7aa2f7"
// positive = "#9ece6a"

cfg := loadConfig("config.toml")
theme := tuikit.ThemeFromMap(cfg.Theme)

app := tuikit.NewApp(
    tuikit.WithTheme(theme),
    ...
)
```

For live hot-reload during development:

```bash
TUIKIT_THEME=./mytheme.json go run .
```

Edit `mytheme.json` and the app reloads the theme without restart.

---

## 5. SSH-Served TUI

Serve your tuikit-go app over SSH using [Wish](https://github.com/charmbracelet/wish) so remote users can access it without installing anything:

```go
package main

import (
    "context"
    "os"
    "os/signal"
    "syscall"

    "github.com/charmbracelet/ssh"
    "github.com/charmbracelet/wish"
    "github.com/charmbracelet/wish/bubbletea"
    tuikit "github.com/moneycaringcoder/tuikit-go"
)

func main() {
    s, _ := wish.NewServer(
        wish.WithAddress(":2222"),
        wish.WithHostKeyPath(".ssh/id_ed25519"),
        wish.WithMiddleware(
            bubbletea.Middleware(func(s ssh.Session) (tea.Model, []tea.ProgramOption) {
                table := tuikit.NewTable(columns, rows, tuikit.TableOpts{Sortable: true})
                app := tuikit.NewApp(
                    tuikit.WithTheme(tuikit.DefaultTheme()),
                    tuikit.WithComponent("main", table),
                )
                return app.Model(), []tea.ProgramOption{tea.WithAltScreen()}
            }),
        ),
    )

    done := make(chan os.Signal, 1)
    signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
    go s.ListenAndServe()
    <-done
    s.Shutdown(context.Background())
}
```

Connect:

```bash
ssh localhost -p 2222
```
