package tuitest

import (
	"fmt"
	"strings"
	"testing"
)

// Diff represents a line-by-line comparison between two screen states.
type Diff struct {
	Lines []DiffLine
}

// DiffLine represents a single line difference.
type DiffLine struct {
	Row      int
	Expected string
	Actual   string
	Changed  bool
}

// ScreenDiff compares two screens and returns a Diff showing changed lines.
func ScreenDiff(a, b *Screen) Diff {
	_, aLines := a.Size()
	_, bLines := b.Size()
	maxLines := aLines
	if bLines > maxLines {
		maxLines = bLines
	}

	var lines []DiffLine
	for r := 0; r < maxLines; r++ {
		var aRow, bRow string
		if r < aLines {
			aRow = a.Row(r)
		}
		if r < bLines {
			bRow = b.Row(r)
		}
		lines = append(lines, DiffLine{
			Row:      r,
			Expected: aRow,
			Actual:   bRow,
			Changed:  aRow != bRow,
		})
	}
	return Diff{Lines: lines}
}

// HasChanges reports whether any lines differ.
func (d Diff) HasChanges() bool {
	for _, l := range d.Lines {
		if l.Changed {
			return true
		}
	}
	return false
}

// ChangedLines returns only the lines that differ.
func (d Diff) ChangedLines() []DiffLine {
	var changed []DiffLine
	for _, l := range d.Lines {
		if l.Changed {
			changed = append(changed, l)
		}
	}
	return changed
}

// String returns a human-readable unified-style diff.
func (d Diff) String() string {
	var b strings.Builder
	for _, l := range d.Lines {
		if !l.Changed {
			continue
		}
		fmt.Fprintf(&b, "  row %d:\n", l.Row)
		fmt.Fprintf(&b, "    - %q\n", l.Expected)
		fmt.Fprintf(&b, "    + %q\n", l.Actual)
	}
	return b.String()
}

// diffStrings produces a simple line-by-line diff between two multi-line strings.
// Used internally by assertions for readable error output.
func diffStrings(expected, actual string) string {
	eLines := strings.Split(expected, "\n")
	aLines := strings.Split(actual, "\n")
	maxLen := len(eLines)
	if len(aLines) > maxLen {
		maxLen = len(aLines)
	}

	var b strings.Builder
	for i := 0; i < maxLen; i++ {
		var e, a string
		if i < len(eLines) {
			e = eLines[i]
		}
		if i < len(aLines) {
			a = aLines[i]
		}
		if e != a {
			fmt.Fprintf(&b, "  row %d:\n", i)
			fmt.Fprintf(&b, "    - %q\n", e)
			fmt.Fprintf(&b, "    + %q\n", a)
		}
	}
	return b.String()
}

// AssertScreensEqual fails if two screens don't have identical text content.
func AssertScreensEqual(t testing.TB, a, b *Screen) {
	t.Helper()
	diff := ScreenDiff(a, b)
	if diff.HasChanges() {
		t.Errorf("screens differ:\n%s", diff.String())
	}
}

// AssertScreensNotEqual fails if two screens have identical text content.
func AssertScreensNotEqual(t testing.TB, a, b *Screen) {
	t.Helper()
	diff := ScreenDiff(a, b)
	if !diff.HasChanges() {
		t.Error("expected screens to differ, but they are identical")
	}
}
