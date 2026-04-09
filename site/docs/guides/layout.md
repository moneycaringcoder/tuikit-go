# Layout

## Single Pane

Register one component as the main pane:

```go
tuikit.WithComponent("main", myComponent)
```

## Dual Pane

Side-by-side layout with a collapsible sidebar:

```go
tuikit.WithLayout(&tuikit.DualPane{
    Main:         table,
    Side:         panel,
    SideWidth:    30,        // character columns
    MinMainWidth: 60,        // sidebar auto-hides below this terminal width
    SideRight:    true,      // sidebar on the right (false = left)
    ToggleKey:    "p",       // key to collapse/expand
})
```

`DualPane.Main` maps to `SlotMain`; `DualPane.Side` maps to `SlotSidebar`. Focus cycles between the two panes with `Tab`.

## StatusBar

Attach a footer with left and right text sections:

```go
tuikit.WithStatusBar(
    func() string { return " ? help  q quit" },
    func() string { return fmt.Sprintf(" %d rows", count) },
)
```

For reactive content driven by signals (e.g. background polling):

```go
leftSig  := tuikit.NewSignal("")
rightSig := tuikit.NewSignal("")

tuikit.WithStatusBarSignal(leftSig, rightSig)

// From any goroutine:
leftSig.Set("connected")
```

Signal updates are coalesced into one notification per frame via a dirty-bit mechanism.

## Tick Interval

Register a periodic tick for polling or animation:

```go
tuikit.WithTickInterval(100 * time.Millisecond)
```

Components receive `tuikit.TickMsg` in their `Update` method.
