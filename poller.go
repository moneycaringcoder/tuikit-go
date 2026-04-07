package tuikit

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Poller manages periodic command execution on top of tuikit's TickMsg.
// Create one with NewPoller, then call Check from your component's Update
// when it receives a TickMsg. The Poller handles interval timing, pause/resume,
// and force-refresh.
type Poller struct {
	interval     time.Duration
	cmd          func() tea.Cmd
	lastPoll     time.Time
	paused       bool
	needsRefresh bool
}

// NewPoller creates a Poller that runs cmd at the given interval.
// The cmd function is called when it's time to poll — it should return
// a tea.Cmd that fetches data (e.g., an API call wrapped in a command).
func NewPoller(interval time.Duration, cmd func() tea.Cmd) *Poller {
	return &Poller{
		interval: interval,
		cmd:      cmd,
	}
}

// Check should be called from your component's Update when receiving a TickMsg.
// Returns a tea.Cmd if it's time to poll, nil otherwise.
// ForceRefresh takes priority and works even when paused.
func (p *Poller) Check(msg TickMsg) tea.Cmd {
	if p.needsRefresh {
		p.needsRefresh = false
		p.lastPoll = time.Now()
		return p.cmd()
	}

	if p.paused {
		return nil
	}

	if time.Since(p.lastPoll) >= p.interval {
		p.lastPoll = time.Now()
		return p.cmd()
	}
	return nil
}

// SetInterval changes the polling interval.
func (p *Poller) SetInterval(d time.Duration) {
	p.interval = d
}

// Pause stops periodic polling. ForceRefresh still works when paused.
func (p *Poller) Pause() { p.paused = true }

// Resume resumes periodic polling.
func (p *Poller) Resume() { p.paused = false }

// TogglePause toggles between paused and active.
func (p *Poller) TogglePause() { p.paused = !p.paused }

// ForceRefresh triggers a poll on the next Check call, even if paused.
func (p *Poller) ForceRefresh() { p.needsRefresh = true }

// IsPaused returns whether polling is paused.
func (p *Poller) IsPaused() bool { return p.paused }

// LastPoll returns the time of the last successful poll.
func (p *Poller) LastPoll() time.Time { return p.lastPoll }
