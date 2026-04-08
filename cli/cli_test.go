package cli

import (
	"errors"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// keyMsg builds a tea.KeyMsg for a named key string.
func keyMsg(key string) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
}

// namedKey builds a tea.KeyMsg for special named keys like "enter", "up", "down".
func namedKey(typ tea.KeyType) tea.KeyMsg {
	return tea.KeyMsg{Type: typ}
}

// --- confirm tests ---

func TestConfirmHint(t *testing.T) {
	tests := []struct {
		defaultYes bool
		wantHint   string
	}{
		{true, "[Y/n]"},
		{false, "[y/N]"},
	}
	for _, tt := range tests {
		m := confirmModel{prompt: "Continue?", defaultYes: tt.defaultYes}
		view := m.View()
		if !strings.Contains(view, tt.wantHint) {
			t.Errorf("defaultYes=%v: expected hint %q in view %q", tt.defaultYes, tt.wantHint, view)
		}
	}
}

func TestConfirmKeyY(t *testing.T) {
	m := confirmModel{prompt: "q?", defaultYes: false}
	result, _ := m.Update(keyMsg("y"))
	final := result.(confirmModel)
	if !final.result {
		t.Error("expected result=true for 'y'")
	}
	if !final.done {
		t.Error("expected done=true")
	}
}

func TestConfirmKeyN(t *testing.T) {
	m := confirmModel{prompt: "q?", defaultYes: true}
	result, _ := m.Update(keyMsg("n"))
	final := result.(confirmModel)
	if final.result {
		t.Error("expected result=false for 'n'")
	}
	if !final.done {
		t.Error("expected done=true")
	}
}

func TestConfirmEnterDefaultYes(t *testing.T) {
	m := confirmModel{prompt: "q?", defaultYes: true}
	result, _ := m.Update(namedKey(tea.KeyEnter))
	final := result.(confirmModel)
	if !final.result {
		t.Error("expected result=true for Enter with defaultYes=true")
	}
}

func TestConfirmEnterDefaultNo(t *testing.T) {
	m := confirmModel{prompt: "q?", defaultYes: false}
	result, _ := m.Update(namedKey(tea.KeyEnter))
	final := result.(confirmModel)
	if final.result {
		t.Error("expected result=false for Enter with defaultYes=false")
	}
}

func TestConfirmViewShowsResultWhenDone(t *testing.T) {
	m := confirmModel{prompt: "Continue?", done: true, result: true}
	view := m.View()
	if !strings.Contains(view, "Continue?") || !strings.Contains(view, "Yes") {
		t.Errorf("expected completion view with prompt and answer, got: %q", view)
	}
	m.result = false
	view = m.View()
	if !strings.Contains(view, "No") {
		t.Errorf("expected 'No' in completion view, got: %q", view)
	}
}

// --- input validation tests ---

func TestInputValidation(t *testing.T) {
	notEmpty := func(s string) error {
		if s == "" {
			return errors.New("required")
		}
		return nil
	}

	if err := notEmpty("hello"); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if err := notEmpty(""); err == nil {
		t.Error("expected error for empty input")
	}
}

func TestInputViewShowsError(t *testing.T) {
	m := newInputModel("Name:", nil)
	m.errMsg = "required"
	view := m.View()
	if !strings.Contains(view, "required") {
		t.Errorf("expected error message in view, got: %q", view)
	}
}

func TestInputViewShowsValueWhenDone(t *testing.T) {
	m := newInputModel("Name:", nil)
	m.done = true
	m.value = "myproject"
	view := m.View()
	if !strings.Contains(view, "Name:") || !strings.Contains(view, "myproject") {
		t.Errorf("expected completion view with prompt and value, got: %q", view)
	}
}

func TestInputViewEmptyWhenCancelled(t *testing.T) {
	m := newInputModel("Name:", nil)
	m.done = true
	m.quitting = true
	if m.View() != "" {
		t.Errorf("expected empty view when cancelled")
	}
}

// --- select tests ---

func TestSelectFilter(t *testing.T) {
	items := make([]string, 15)
	for i := range items {
		items[i] = "item"
	}
	items[3] = "special"
	items[12] = "special-two"

	m := newSelectModel("Pick:", items)
	if !m.useFilter {
		t.Fatal("expected filter to be active for 15 items")
	}

	m.filter = "special"
	m.applyFilter()

	if len(m.filtered) != 2 {
		t.Errorf("expected 2 filtered items, got %d", len(m.filtered))
	}
}

func TestSelectFilterInactiveSmallList(t *testing.T) {
	items := []string{"a", "b", "c"}
	m := newSelectModel("Pick:", items)
	if m.useFilter {
		t.Error("expected filter to be inactive for 3 items")
	}
}

func TestSelectCursorDown(t *testing.T) {
	items := []string{"a", "b", "c"}
	m := newSelectModel("Pick:", items)

	result, _ := m.Update(namedKey(tea.KeyDown))
	m = result.(selectModel)
	if m.cursor != 1 {
		t.Errorf("expected cursor=1 after down, got %d", m.cursor)
	}
}

func TestSelectCursorUp(t *testing.T) {
	items := []string{"a", "b", "c"}
	m := newSelectModel("Pick:", items)
	m.cursor = 2

	result, _ := m.Update(namedKey(tea.KeyUp))
	m = result.(selectModel)
	if m.cursor != 1 {
		t.Errorf("expected cursor=1 after up, got %d", m.cursor)
	}
}

func TestSelectCursorClampTop(t *testing.T) {
	items := []string{"a", "b"}
	m := newSelectModel("Pick:", items)

	result, _ := m.Update(namedKey(tea.KeyUp))
	m = result.(selectModel)
	if m.cursor != 0 {
		t.Errorf("cursor should clamp at 0, got %d", m.cursor)
	}
}

func TestSelectCursorClampBottom(t *testing.T) {
	items := []string{"a", "b"}
	m := newSelectModel("Pick:", items)
	m.cursor = 1

	result, _ := m.Update(namedKey(tea.KeyDown))
	m = result.(selectModel)
	if m.cursor != 1 {
		t.Errorf("cursor should clamp at 1, got %d", m.cursor)
	}
}

func TestSelectViewContainsItems(t *testing.T) {
	items := []string{"apple", "banana", "cherry"}
	m := newSelectModel("Pick:", items)
	view := m.View()
	for _, item := range items {
		if !strings.Contains(view, item) {
			t.Errorf("expected %q in view", item)
		}
	}
}

func TestSelectEmptyItems(t *testing.T) {
	_, _, err := SelectOne("Pick:", []string{})
	if err == nil {
		t.Error("expected error for empty items")
	}
}

// --- progress bar tests ---

func TestProgressPercentage(t *testing.T) {
	p := &Progress{total: 100, label: "Test"}
	p.current = 50

	pct := p.current * 100 / p.total
	if pct != 50 {
		t.Errorf("expected 50%%, got %d%%", pct)
	}
}

func TestProgressBarFilled(t *testing.T) {
	p := &Progress{total: 100, label: "Test"}
	p.current = 50

	filled := p.current * barWidth / p.total
	if filled != barWidth/2 {
		t.Errorf("expected %d filled cells, got %d", barWidth/2, filled)
	}
}

func TestProgressIncrement(t *testing.T) {
	p := &Progress{total: 10, label: "Work"}
	p.Increment(3)
	if p.current != 3 {
		t.Errorf("expected current=3, got %d", p.current)
	}
}

func TestProgressIncrementClamps(t *testing.T) {
	p := &Progress{total: 10, label: "Work"}
	p.Increment(3)
	p.Increment(10) // should clamp to total
	if p.current != 10 {
		t.Errorf("expected current clamped to 10, got %d", p.current)
	}
}

func TestProgressDoneIdempotent(t *testing.T) {
	p := &Progress{total: 5, label: "X"}
	p.Done()
	p.Done() // second call must not panic
	if !p.done {
		t.Error("expected done=true")
	}
}

func TestProgressZeroTotal(t *testing.T) {
	// Should not divide by zero
	p := &Progress{total: 0, label: "Empty"}
	p.Increment(1)
	p.Done()
}

// --- spinner tests ---

func TestSpinnerStopDoesNotPanic(t *testing.T) {
	s := Spin("loading...")
	s.Stop()
	s.Stop() // idempotent — second stop must not panic
}

// --- multiselect tests ---

func TestMultiSelectToggle(t *testing.T) {
	m := newMultiSelectModel("Pick:", []string{"a", "b", "c"})
	// Toggle first item
	result, _ := m.Update(keyMsg(" "))
	m = result.(multiSelectModel)
	if !m.checked[0] {
		t.Error("expected item 0 to be checked after space")
	}
	// Toggle again to uncheck
	result, _ = m.Update(keyMsg(" "))
	m = result.(multiSelectModel)
	if m.checked[0] {
		t.Error("expected item 0 to be unchecked after second space")
	}
}

func TestMultiSelectToggleX(t *testing.T) {
	m := newMultiSelectModel("Pick:", []string{"a", "b"})
	result, _ := m.Update(keyMsg("x"))
	m = result.(multiSelectModel)
	if !m.checked[0] {
		t.Error("expected item 0 to be checked after 'x'")
	}
}

func TestMultiSelectToggleAll(t *testing.T) {
	m := newMultiSelectModel("Pick:", []string{"a", "b", "c"})
	result, _ := m.Update(keyMsg("a"))
	m = result.(multiSelectModel)
	for i, c := range m.checked {
		if !c {
			t.Errorf("expected item %d to be checked after 'a' (toggle all)", i)
		}
	}
	// Toggle all again to uncheck
	result, _ = m.Update(keyMsg("a"))
	m = result.(multiSelectModel)
	for i, c := range m.checked {
		if c {
			t.Errorf("expected item %d to be unchecked after second 'a'", i)
		}
	}
}

func TestMultiSelectCursorNav(t *testing.T) {
	m := newMultiSelectModel("Pick:", []string{"a", "b", "c"})
	result, _ := m.Update(namedKey(tea.KeyDown))
	m = result.(multiSelectModel)
	if m.cursor != 1 {
		t.Errorf("cursor = %d, want 1", m.cursor)
	}
	result, _ = m.Update(namedKey(tea.KeyDown))
	m = result.(multiSelectModel)
	if m.cursor != 2 {
		t.Errorf("cursor = %d, want 2", m.cursor)
	}
	// Clamp at bottom
	result, _ = m.Update(namedKey(tea.KeyDown))
	m = result.(multiSelectModel)
	if m.cursor != 2 {
		t.Errorf("cursor should clamp at 2, got %d", m.cursor)
	}
}

func TestMultiSelectSelected(t *testing.T) {
	m := newMultiSelectModel("Pick:", []string{"a", "b", "c"})
	m.checked[0] = true
	m.checked[2] = true

	got := m.selected()
	if len(got) != 2 || got[0] != "a" || got[1] != "c" {
		t.Errorf("selected() = %v, want [a c]", got)
	}
	indices := m.selectedIndices()
	if len(indices) != 2 || indices[0] != 0 || indices[1] != 2 {
		t.Errorf("selectedIndices() = %v, want [0 2]", indices)
	}
}

func TestMultiSelectViewContainsItems(t *testing.T) {
	m := newMultiSelectModel("Choose:", []string{"alpha", "beta"})
	view := m.View()
	if !strings.Contains(view, "alpha") || !strings.Contains(view, "beta") {
		t.Errorf("expected items in view, got: %q", view)
	}
	if !strings.Contains(view, "Choose:") {
		t.Error("expected prompt in view")
	}
}

func TestMultiSelectViewDoneShowsSelected(t *testing.T) {
	m := newMultiSelectModel("Pick:", []string{"a", "b", "c"})
	m.checked[1] = true
	m.done = true
	view := m.View()
	if !strings.Contains(view, "b") {
		t.Errorf("expected selected item in done view, got: %q", view)
	}
}

func TestMultiSelectViewDoneNone(t *testing.T) {
	m := newMultiSelectModel("Pick:", []string{"a", "b"})
	m.done = true
	view := m.View()
	if !strings.Contains(view, "(none)") {
		t.Errorf("expected '(none)' when nothing selected, got: %q", view)
	}
}

func TestMultiSelectEmptyItems(t *testing.T) {
	_, _, err := MultiSelect("Pick:", []string{})
	if err == nil {
		t.Error("expected error for empty items")
	}
}

// --- password tests ---

func TestPasswordMasked(t *testing.T) {
	m := newPasswordModel("Secret:", nil)
	if m.ti.EchoCharacter != '•' {
		t.Errorf("EchoCharacter = %c, want •", m.ti.EchoCharacter)
	}
}

func TestPasswordViewShowsMasked(t *testing.T) {
	m := newPasswordModel("Secret:", nil)
	m.done = true
	m.value = "hello"
	view := m.View()
	if !strings.Contains(view, "•••••") {
		t.Errorf("expected masked dots in done view, got: %q", view)
	}
	if strings.Contains(view, "hello") {
		t.Error("password value should not appear in plain text")
	}
}

func TestPasswordViewEmptyWhenCancelled(t *testing.T) {
	m := newPasswordModel("Secret:", nil)
	m.done = true
	m.quitting = true
	if m.View() != "" {
		t.Error("expected empty view when cancelled")
	}
}

func TestPasswordValidation(t *testing.T) {
	minLen := func(s string) error {
		if len(s) < 8 {
			return errors.New("too short")
		}
		return nil
	}
	m := newPasswordModel("Password:", minLen)
	// Simulate entering a short password
	m.ti.SetValue("abc")
	result, _ := m.Update(namedKey(tea.KeyEnter))
	final := result.(passwordModel)
	if final.done {
		t.Error("should not be done with invalid password")
	}
	if final.errMsg != "too short" {
		t.Errorf("errMsg = %q, want %q", final.errMsg, "too short")
	}
}

func TestPasswordViewShowsError(t *testing.T) {
	m := newPasswordModel("Pass:", nil)
	m.errMsg = "too short"
	view := m.View()
	if !strings.Contains(view, "too short") {
		t.Errorf("expected error in view, got: %q", view)
	}
}
