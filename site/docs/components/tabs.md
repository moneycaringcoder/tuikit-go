# Tabs

A component that shows a row or column of named tabs with switchable content panes. Supports mouse click-to-select, keyboard cycling, and direct index jump. Implements `Component` and `Themed`.

## Construction

```go
tabs := tuikit.NewTabs(items []tuikit.TabItem, opts tuikit.TabsOpts)
```

## TabItem

```go
type TabItem struct {
    Title   string           // Label shown in the tab bar
    Glyph   string           // Optional icon prefix (e.g. "⚡")
    Content tuikit.Component // Component rendered when this tab is active
}
```

## TabsOpts

```go
type TabsOpts struct {
    Orientation tuikit.Orientation // Horizontal (default) or Vertical
    OnChange    func(int)          // Called whenever the active tab changes
}
```

## Horizontal Tabs (Default)

```go
tabs := tuikit.NewTabs([]tuikit.TabItem{
    {Title: "Overview",  Glyph: "◉", Content: overviewPanel},
    {Title: "Logs",      Glyph: "≡", Content: logViewer},
    {Title: "Settings",  Glyph: "⚙", Content: configEditor},
}, tuikit.TabsOpts{
    OnChange: func(idx int) {
        fmt.Println("switched to tab", idx)
    },
})
```

The horizontal bar renders three rows: the tab labels row, an accent underline row, and a full-width border line.

## Vertical Tabs (Sidebar Style)

```go
tabs := tuikit.NewTabs(items, tuikit.TabsOpts{
    Orientation: tuikit.Vertical,
})
```

In vertical mode, a sidebar column shows tab labels. The active tab label has accent background + bold text. The sidebar width is computed from the longest label (minimum 12 columns).

## Programmatic Control

```go
tabs.SetActive(2)            // Switch to tab index 2 (clamped to valid range)
idx := tabs.ActiveIndex()    // Get current active tab index
tabs.OnChange(func(i int) {  // Register or replace OnChange callback
    updateBreadcrumb(i)
})
```

## Focus Propagation

When `SetFocused(true)` is called on `Tabs`, only the active content pane receives focus. Switching tabs automatically transfers focus to the newly active content.

## Theme Propagation

`SetTheme` is forwarded to all content components that implement the `Themed` interface.

## Keybindings

| Key | Action |
|-----|--------|
| `tab` | Switch to next tab |
| `shift+tab` | Switch to previous tab |
| `1`–`9` | Jump directly to tab N |
| Left-click | Click tab label to select |

## Embedding in a Layout

```go
app := tuikit.NewApp(
    tuikit.WithComponent("main", tabs),
    tuikit.WithStatusBar(
        func() string { return " tab next  shift+tab prev" },
        func() string { return fmt.Sprintf("tab %d/3", tabs.ActiveIndex()+1) },
    ),
)
```
