# API Reference

Full API documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/moneycaringcoder/tuikit-go).

## Core

| Type | Description |
|------|-------------|
| [`App`](https://pkg.go.dev/github.com/moneycaringcoder/tuikit-go#App) | Application container configured via functional options |
| [`Component`](https://pkg.go.dev/github.com/moneycaringcoder/tuikit-go#Component) | Interface all components implement (Init, Update, View, KeyBindings, SetSize, Focused, SetFocused) |
| [`Context`](https://pkg.go.dev/github.com/moneycaringcoder/tuikit-go#Context) | Per-update context carrying Theme, Size, Focus, Hotkeys, Clock, Logger |
| [`Theme`](https://pkg.go.dev/github.com/moneycaringcoder/tuikit-go#Theme) | Semantic color tokens (Positive, Negative, Accent, Muted, etc.) |

## Components

| Type | Description |
|------|-------------|
| [`Table`](https://pkg.go.dev/github.com/moneycaringcoder/tuikit-go#Table) | Adaptive table with sorting, filtering, custom rendering, virtualization |
| [`ListView`](https://pkg.go.dev/github.com/moneycaringcoder/tuikit-go#ListView) | Vertical list with cursor navigation |
| [`Tabs`](https://pkg.go.dev/github.com/moneycaringcoder/tuikit-go#Tabs) | Tabbed container with horizontal/vertical orientation |
| [`Form`](https://pkg.go.dev/github.com/moneycaringcoder/tuikit-go#Form) | Form with validation and wizard mode |
| [`Picker`](https://pkg.go.dev/github.com/moneycaringcoder/tuikit-go#Picker) | Fuzzy-search selection list |
| [`Tree`](https://pkg.go.dev/github.com/moneycaringcoder/tuikit-go#Tree) | Expandable tree view |
| [`Viewport`](https://pkg.go.dev/github.com/moneycaringcoder/tuikit-go#Viewport) | Scrollable content pane |
| [`LogViewer`](https://pkg.go.dev/github.com/moneycaringcoder/tuikit-go#LogViewer) | Streaming log viewer with level filtering |
| [`Markdown`](https://pkg.go.dev/github.com/moneycaringcoder/tuikit-go#Markdown) | Glamour-powered markdown renderer |

## Layout

| Type | Description |
|------|-------------|
| [`DualPane`](https://pkg.go.dev/github.com/moneycaringcoder/tuikit-go#DualPane) | Main + collapsible sidebar layout |
| [`HBox`](https://pkg.go.dev/github.com/moneycaringcoder/tuikit-go#HBox) | Horizontal flex container |
| [`VBox`](https://pkg.go.dev/github.com/moneycaringcoder/tuikit-go#VBox) | Vertical flex container |
| [`Split`](https://pkg.go.dev/github.com/moneycaringcoder/tuikit-go#Split) | Resizable split pane |

## Packages

| Package | Description |
|---------|-------------|
| [`cli`](https://pkg.go.dev/github.com/moneycaringcoder/tuikit-go/cli) | Interactive CLI prompts (Confirm, Select, Input, Spinner, Progress) |
| [`charts`](https://pkg.go.dev/github.com/moneycaringcoder/tuikit-go/charts) | Chart components (Bar, Line, Ring, Gauge, Heatmap) |
| [`tuitest`](https://pkg.go.dev/github.com/moneycaringcoder/tuikit-go/tuitest) | Virtual terminal testing framework |
