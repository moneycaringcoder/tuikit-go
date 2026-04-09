# Testing with tuitest

`tuitest` is a virtual-terminal testing framework for Bubble Tea models. It renders your model into an in-memory screen buffer and provides 30+ assertions — no real terminal required.

## Install the tuitest CLI

```bash
# Homebrew
brew install moneycaringcoder/tap/tuitest

# Scoop
scoop bucket add moneycaringcoder https://github.com/moneycaringcoder/scoop-bucket
scoop install tuitest

# Go
go install github.com/moneycaringcoder/tuikit-go/cmd/tuitest@latest
```

## Writing a Test

```go
import "github.com/moneycaringcoder/tuikit-go/tuitest"

func TestMyApp(t *testing.T) {
    tm := tuitest.NewTestModel(t, myModel{}, 80, 24)

    // Interact
    tm.SendKey("down")
    tm.SendKeys("j", "j", "enter")
    tm.Type("hello")
    tm.SendResize(120, 40)
    tm.SendMsg(myCustomMsg{})

    // Assert on rendered screen
    scr := tm.Screen()
    tuitest.AssertContains(t, scr, "Expected text")
    tuitest.AssertRowContains(t, scr, 0, "Header")
    tuitest.AssertMatches(t, scr, `\d+ items`)
    tuitest.AssertRowCount(t, scr, 5)
}
```

## Assertions Reference

| Assertion | Description |
|-----------|-------------|
| `AssertContains` | Screen contains substring |
| `AssertNotContains` | Screen does not contain substring |
| `AssertRowContains` | Row N contains substring |
| `AssertMatches` | Screen matches regexp |
| `AssertRowCount` | Screen has exactly N non-empty rows |
| `AssertFgAt` | Cell at (row, col) has foreground color |
| `AssertBgAt` | Cell at (row, col) has background color |
| `AssertBoldAt` | Cell at (row, col) is bold |
| `AssertRegionContains` | Rectangular region contains text |
| `AssertScreensEqual` | Two screen snapshots are identical |
| `AssertScreensNotEqual` | Two screen snapshots differ |
| `AssertGolden` | Compare against golden file in `testdata/` |

## Golden File Testing

```go
tuitest.AssertGolden(t, scr, "my-test")
// compares against testdata/my-test.golden
```

Regenerate snapshots:

```bash
tuitest -update ./...
```

## Waiting for Async State

```go
ok := tm.WaitFor(tuitest.UntilContains("loaded"), 10)
if !ok {
    t.Fatal("timed out waiting for 'loaded'")
}
```

## tuitest CLI Flags

```bash
tuitest                                    # go test ./...
tuitest -filter TestHarness ./tuitest/...  # run tests matching a regexp
tuitest -update ./tuitest/...              # regenerate golden snapshots
tuitest -junit out/junit.xml -parallel 4   # parallel run + JUnit report
tuitest -html out/report.html              # HTML report
tuitest -watch                             # re-run on file changes (1s poll)
```

## Vitest-Style Reporter

Run with `-v` for grouped, color-coded output:

```
  tuitest · terminal test toolkit

  Screen
    ✓ PlainText 0.000ms
    ✓ Contains 0.000ms
  Assert
    ✓ ContainsPass 0.000ms
    ✓ RowMatchesPass 0.000ms

  PASS 96 tests (3ms)
```
